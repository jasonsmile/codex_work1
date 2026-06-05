package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"drug-info/backend/config"
	"drug-info/backend/logger"
	"drug-info/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const maxTraceCodeUploadSize = 10 << 20
const traceCodeOCRTimeout = 60 * time.Second

var errBaiduOCRConfigMissing = errors.New("missing_baidu_ocr_config")

var allowedTraceCodeUploadTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

type TraceCodeHandler struct {
	db       *gorm.DB
	baiduOCR config.BaiduOCRConfig
	token    baiduAccessTokenCache
}

type baiduAccessTokenCache struct {
	mu        sync.Mutex
	value     string
	expiresAt time.Time
}

type traceCodeRecognizeResponse struct {
	Records []traceCodeRecord `json:"records"`
}

type traceCodeRecord struct {
	TransactionSerialNumber string `json:"transaction_serial_number"`
	DrugCode                string `json:"drug_code"`
	DrugName                string `json:"drug_name"`
	SettlementDate          string `json:"settlement_date"`
	SourceFileName          string `json:"source_file_name"`
	ContentType             string `json:"content_type"`
}

type traceCodeConfirmRequest struct {
	Records []traceCodeRecord `json:"records"`
}

type traceCodeDeleteRequest struct {
	ID uint `json:"id"`
}

type traceCodeListResponse struct {
	List     []models.TraceCode `json:"list"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

type baiduTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

type baiduOCRResponse struct {
	WordsResult []baiduOCRWord `json:"words_result"`
	ErrorCode   int64          `json:"error_code"`
	ErrorMsg    string         `json:"error_msg"`
}

type baiduOCRWord struct {
	Words    string           `json:"words"`
	Location baiduOCRLocation `json:"location"`
}

type baiduOCRLocation struct {
	Top    int `json:"top"`
	Left   int `json:"left"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func NewTraceCodeHandler(database *gorm.DB, baiduOCRConfig config.BaiduOCRConfig) *TraceCodeHandler {
	return &TraceCodeHandler{db: database, baiduOCR: baiduOCRConfig}
}

func (h *TraceCodeHandler) Recognize(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxTraceCodeUploadSize+1024*1024)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		badRequest(c, "请选择要识别的图片文件", err)
		return
	}
	if fileHeader.Size > maxTraceCodeUploadSize {
		badRequest(c, "文件大小不能超过 10MB", nil)
		return
	}

	fileName := strings.TrimSpace(fileHeader.Filename)
	ext := strings.ToLower(filepath.Ext(fileName))
	contentType, ok := allowedTraceCodeUploadTypes[ext]
	if !ok {
		badRequest(c, "仅支持 jpg、jpeg、png、webp 图片", nil)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		serverError(c, "读取上传文件失败", err)
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		serverError(c, "读取上传文件失败", err)
		return
	}

	result, err := h.recognizeTraceCode(c.Request.Context(), imageBytes)
	if err != nil {
		if errors.Is(err, errBaiduOCRConfigMissing) {
			badRequest(c, "未配置百度 OCR AK/SK，请在 config.yaml 的 baidu-ocr.ak 和 baidu-ocr.sk 中配置", nil)
			return
		}
		serverError(c, "识别失败", err)
		return
	}
	for index := range result.Records {
		result.Records[index].SourceFileName = fileName
		result.Records[index].ContentType = contentType
	}
	if len(result.Records) == 0 {
		logger.Warning("trace code recognize no records",
			logger.Field{Key: "request_id", Value: logger.RequestID(c)},
			logger.Field{Key: "file_name", Value: fileName},
			logger.Field{Key: "file_size", Value: fileHeader.Size},
			logger.Field{Key: "content_type", Value: contentType},
			logger.Field{Key: "client_ip", Value: c.ClientIP()},
		)
		badRequest(c, "未识别到可入库记录", nil)
		return
	}

	success(c, http.StatusOK, "识别成功", result)
}

func (h *TraceCodeHandler) Confirm(c *gin.Context) {
	var request traceCodeConfirmRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		badRequest(c, "请求参数不正确", err)
		return
	}
	if len(request.Records) == 0 {
		badRequest(c, "暂无可入库的识别结果", nil)
		return
	}

	traceCodes := make([]models.TraceCode, 0, len(request.Records))
	for _, record := range request.Records {
		traceCode, err := buildTraceCode(record)
		if err != nil {
			badRequest(c, err.Error(), nil)
			return
		}
		traceCodes = append(traceCodes, traceCode)
	}

	if err := h.db.Create(&traceCodes).Error; err != nil {
		serverError(c, "保存异常追溯码失败", err)
		return
	}

	success(c, http.StatusCreated, "入库成功", traceCodes)
}

func (h *TraceCodeHandler) List(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	pageSize := parsePositiveInt(c.Query("page_size"), 20)
	if pageSize > 100 {
		pageSize = 100
	}

	query := h.db.Model(&models.TraceCode{})
	transactionSerialNumber := strings.TrimSpace(c.Query("transaction_serial_number"))
	if transactionSerialNumber != "" {
		query = query.Where("transaction_serial_number LIKE ?", "%"+transactionSerialNumber+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		serverError(c, "查询异常追溯码总数失败", err)
		return
	}

	var traceCodes []models.TraceCode
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&traceCodes).Error; err != nil {
		serverError(c, "查询异常追溯码列表失败", err)
		return
	}

	success(c, http.StatusOK, "查询成功", traceCodeListResponse{
		List:     traceCodes,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (h *TraceCodeHandler) Delete(c *gin.Context) {
	var request traceCodeDeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		badRequest(c, "请求参数不正确", err)
		return
	}
	if request.ID == 0 {
		badRequest(c, "记录ID不能为空", nil)
		return
	}

	result := h.db.Delete(&models.TraceCode{}, request.ID)
	if result.Error != nil {
		serverError(c, "删除异常追溯码失败", result.Error)
		return
	}
	if result.RowsAffected == 0 {
		badRequest(c, "记录不存在", nil)
		return
	}

	success(c, http.StatusOK, "删除成功", nil)
}

func (h *TraceCodeHandler) recognizeTraceCode(ctx context.Context, imageBytes []byte) (traceCodeRecognizeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, traceCodeOCRTimeout)
	defer cancel()

	words, err := h.callBaiduOCR(ctx, imageBytes)
	if err != nil {
		return traceCodeRecognizeResponse{}, err
	}
	records := parseTraceCodeWords(words)
	cleanTraceCodeRecords(records)
	return traceCodeRecognizeResponse{Records: records}, nil
}

func (h *TraceCodeHandler) callBaiduOCR(ctx context.Context, imageBytes []byte) ([]baiduOCRWord, error) {
	accessToken, err := h.getBaiduAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	ocrURL := strings.TrimSpace(getEnvOrDefault("BAIDU_OCR_URL", h.baiduOCR.OCRURL))
	if ocrURL == "" {
		ocrURL = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate"
	}

	requestURL := ocrURL + "?access_token=" + url.QueryEscape(accessToken)
	form := url.Values{}
	form.Set("image", base64.StdEncoding.EncodeToString(imageBytes))
	form.Set("language_type", "CHN_ENG")
	form.Set("detect_direction", "true")
	form.Set("paragraph", "false")

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("baidu_ocr_http_status_%d: %s", response.StatusCode, truncateString(string(responseBytes), 500))
	}

	var ocrResponse baiduOCRResponse
	if err := json.Unmarshal(responseBytes, &ocrResponse); err != nil {
		return nil, err
	}
	if ocrResponse.ErrorCode != 0 {
		return nil, fmt.Errorf("baidu_ocr_error_%d: %s", ocrResponse.ErrorCode, ocrResponse.ErrorMsg)
	}

	words := make([]baiduOCRWord, 0, len(ocrResponse.WordsResult))
	for _, item := range ocrResponse.WordsResult {
		word := strings.TrimSpace(item.Words)
		if word != "" {
			item.Words = word
			words = append(words, item)
		}
	}
	return words, nil
}

func (h *TraceCodeHandler) getBaiduAccessToken(ctx context.Context) (string, error) {
	ak := strings.TrimSpace(getEnvOrDefault("BAIDU_OCR_AK", h.baiduOCR.AK))
	sk := strings.TrimSpace(getEnvOrDefault("BAIDU_OCR_SK", h.baiduOCR.SK))
	if ak == "" || sk == "" {
		return "", errBaiduOCRConfigMissing
	}

	h.token.mu.Lock()
	if h.token.value != "" && time.Now().Before(h.token.expiresAt) {
		value := h.token.value
		h.token.mu.Unlock()
		return value, nil
	}
	h.token.mu.Unlock()

	tokenURL := strings.TrimSpace(getEnvOrDefault("BAIDU_OCR_TOKEN_URL", h.baiduOCR.TokenURL))
	if tokenURL == "" {
		tokenURL = "https://aip.baidubce.com/oauth/2.0/token"
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", ak)
	form.Set("client_secret", sk)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("baidu_token_http_status_%d: %s", response.StatusCode, truncateString(string(responseBytes), 500))
	}

	var tokenResponse baiduTokenResponse
	if err := json.Unmarshal(responseBytes, &tokenResponse); err != nil {
		return "", err
	}
	if tokenResponse.AccessToken == "" {
		return "", fmt.Errorf("baidu_token_error_%s: %s", tokenResponse.Error, tokenResponse.Description)
	}

	expiresIn := tokenResponse.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 3600
	}

	h.token.mu.Lock()
	h.token.value = tokenResponse.AccessToken
	h.token.expiresAt = time.Now().Add(time.Duration(expiresIn-300) * time.Second)
	h.token.mu.Unlock()
	return tokenResponse.AccessToken, nil
}

func parseTraceCodeWords(words []baiduOCRWord) []traceCodeRecord {
	records := parseTraceCodeWordsByColumns(words)
	if len(records) > 0 {
		return records
	}

	rows := groupBaiduOCRWordsByRow(words)
	text := strings.Join(rows, " ")
	text = normalizeOCRText(text)

	rowPattern := regexp.MustCompile(`(?i)(\d[a-z0-9]{9,})\s+([a-z0-9]{12,})\s+([\p{Han}a-z0-9（）()·\-]+?)\s+\d+(?:\.\d+)?\s+\d+(?:\.\d+)?\s+(\d{4}-\d{1,2}-\d{1,2})`)
	matches := rowPattern.FindAllStringSubmatch(text, -1)
	records = make([]traceCodeRecord, 0, len(matches))
	for _, matched := range matches {
		if len(matched) < 5 {
			continue
		}
		records = append(records, traceCodeRecord{
			TransactionSerialNumber: strings.TrimSpace(matched[1]),
			DrugCode:                strings.TrimSpace(matched[2]),
			DrugName:                strings.TrimSpace(matched[3]),
			SettlementDate:          normalizeDate(matched[4]),
		})
	}
	if len(records) > 0 {
		return records
	}

	return parseTraceCodeWordsByLine(rows)
}

func parseTraceCodeWordsByColumns(words []baiduOCRWord) []traceCodeRecord {
	rowAnchors := make([]baiduOCRWord, 0)
	for _, word := range words {
		if word.Location.Left > 180 {
			continue
		}
		if isTraceTransactionSerial(word.Words) {
			rowAnchors = append(rowAnchors, word)
		}
	}
	if len(rowAnchors) == 0 {
		return nil
	}

	sort.SliceStable(rowAnchors, func(leftIndex, rightIndex int) bool {
		return rowAnchors[leftIndex].Location.Top < rowAnchors[rightIndex].Location.Top
	})

	records := make([]traceCodeRecord, 0, len(rowAnchors))
	for index, anchor := range rowAnchors {
		rowTop := anchor.Location.Top - 35
		rowBottom := anchor.Location.Top + 35
		if index > 0 && rowTop <= rowAnchors[index-1].Location.Top {
			rowTop = rowAnchors[index-1].Location.Top + 2
		}
		if index+1 < len(rowAnchors) {
			nextTop := rowAnchors[index+1].Location.Top - 2
			if rowBottom > nextTop {
				rowBottom = nextTop
			}
		}

		rowItems := make([]baiduOCRWord, 0)
		for _, word := range words {
			if word.Location.Top >= rowTop && word.Location.Top < rowBottom {
				rowItems = append(rowItems, word)
			}
		}
		sort.SliceStable(rowItems, func(leftIndex, rightIndex int) bool {
			return rowItems[leftIndex].Location.Left < rowItems[rightIndex].Location.Left
		})

		record := traceCodeRecord{
			TransactionSerialNumber: cleanTraceCodeCell(anchor.Words),
			DrugCode:                firstCellInRange(rowItems, 220, 430, anchor.Location.Top),
			DrugName:                firstCellInRange(rowItems, 430, 660, anchor.Location.Top),
			SettlementDate:          normalizeDate(firstCellInRange(rowItems, 880, 1120, anchor.Location.Top)),
		}
		if record.DrugCode == "" || record.DrugName == "" {
			continue
		}
		records = append(records, record)
	}
	return records
}

func firstCellInRange(items []baiduOCRWord, leftMin int, leftMax int, targetTop int) string {
	var best *baiduOCRWord
	bestDistance := 0
	preferredItems := make([]baiduOCRWord, 0)
	for _, item := range items {
		if item.Location.Left >= leftMin && item.Location.Left < leftMax {
			if item.Location.Top >= targetTop-2 {
				preferredItems = append(preferredItems, item)
			}
		}
	}
	searchItems := preferredItems
	if len(searchItems) == 0 {
		searchItems = items
	}
	for _, item := range searchItems {
		if item.Location.Left < leftMin || item.Location.Left >= leftMax {
			continue
		}
		distance := absInt(item.Location.Top - targetTop)
		if best == nil || distance < bestDistance {
			current := item
			best = &current
			bestDistance = distance
		}
	}
	if best != nil {
		return cleanTraceCodeCell(best.Words)
	}
	return ""
}

func isTraceTransactionSerial(value string) bool {
	value = cleanTraceCodeCell(value)
	if len(value) < 10 || len(value) > 28 {
		return false
	}
	if !regexp.MustCompile(`(?i)^\d[a-z0-9]+$`).MatchString(value) {
		return false
	}
	return strings.Contains(strings.ToUpper(value), "Y")
}

func cleanTraceCodeCell(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "　", "")
	return value
}

func groupBaiduOCRWordsByRow(words []baiduOCRWord) []string {
	if len(words) == 0 {
		return nil
	}

	sort.SliceStable(words, func(leftIndex, rightIndex int) bool {
		if words[leftIndex].Location.Top == words[rightIndex].Location.Top {
			return words[leftIndex].Location.Left < words[rightIndex].Location.Left
		}
		return words[leftIndex].Location.Top < words[rightIndex].Location.Top
	})

	type rowGroup struct {
		top   int
		items []baiduOCRWord
	}
	groups := make([]rowGroup, 0)
	for _, word := range words {
		if word.Location.Top == 0 && word.Location.Left == 0 {
			groups = append(groups, rowGroup{top: len(groups) * 100, items: []baiduOCRWord{word}})
			continue
		}

		attached := false
		for index := range groups {
			if absInt(groups[index].top-word.Location.Top) <= 12 {
				groups[index].items = append(groups[index].items, word)
				if word.Location.Top < groups[index].top {
					groups[index].top = word.Location.Top
				}
				attached = true
				break
			}
		}
		if !attached {
			groups = append(groups, rowGroup{top: word.Location.Top, items: []baiduOCRWord{word}})
		}
	}

	rows := make([]string, 0, len(groups))
	for _, group := range groups {
		sort.SliceStable(group.items, func(leftIndex, rightIndex int) bool {
			return group.items[leftIndex].Location.Left < group.items[rightIndex].Location.Left
		})

		cells := make([]string, 0, len(group.items))
		for _, item := range group.items {
			cells = append(cells, strings.TrimSpace(item.Words))
		}
		row := strings.TrimSpace(strings.Join(cells, " "))
		if row != "" {
			rows = append(rows, row)
		}
	}
	return rows
}

func parseTraceCodeWordsByLine(words []string) []traceCodeRecord {
	records := make([]traceCodeRecord, 0)
	for _, word := range words {
		line := normalizeOCRText(word)
		serial := matchFirst(line, `(?i)(\d[a-z0-9]{9,})`)
		drugCode := matchFirst(line, `(?i)\d[a-z0-9]{9,}\s+([a-z0-9]{12,})`)
		if serial == "" || drugCode == "" {
			continue
		}
		settlementDate := normalizeDate(matchFirst(line, `(\d{4}-\d{1,2}-\d{1,2})`))
		drugName := extractDrugName(line, serial, drugCode)
		records = append(records, traceCodeRecord{
			TransactionSerialNumber: serial,
			DrugCode:                drugCode,
			DrugName:                drugName,
			SettlementDate:          settlementDate,
		})
	}
	return records
}

func extractDrugName(line string, serial string, drugCode string) string {
	remaining := strings.TrimSpace(strings.Replace(line, serial, "", 1))
	remaining = strings.TrimSpace(strings.Replace(remaining, drugCode, "", 1))
	name := matchFirst(remaining, `([\p{Han}a-z0-9（）()·\-]+?)\s+\d+(?:\.\d+)?`)
	return strings.TrimSpace(name)
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func buildTraceCode(record traceCodeRecord) (models.TraceCode, error) {
	transactionSerialNumber := strings.TrimSpace(record.TransactionSerialNumber)
	drugCode := strings.TrimSpace(record.DrugCode)
	drugName := strings.TrimSpace(record.DrugName)
	settlementDate := strings.TrimSpace(record.SettlementDate)

	if transactionSerialNumber == "" {
		return models.TraceCode{}, newValidationError("交易流水号不能为空")
	}
	if drugCode == "" {
		return models.TraceCode{}, newValidationError("药品编号不能为空")
	}
	if drugName == "" {
		return models.TraceCode{}, newValidationError("药品名称不能为空")
	}
	if settlementDate != "" {
		if _, err := time.Parse("2006-01-02", settlementDate); err != nil {
			return models.TraceCode{}, newValidationError("结算日期格式必须为 YYYY-MM-DD")
		}
	}

	return models.TraceCode{
		TransactionSerialNumber: transactionSerialNumber,
		DrugCode:                drugCode,
		DrugName:                drugName,
		SettlementDate:          settlementDate,
		SourceFileName:          strings.TrimSpace(record.SourceFileName),
		ContentType:             strings.TrimSpace(record.ContentType),
	}, nil
}

func cleanTraceCodeRecords(records []traceCodeRecord) {
	for index := range records {
		records[index].TransactionSerialNumber = strings.TrimSpace(records[index].TransactionSerialNumber)
		records[index].DrugCode = strings.TrimSpace(records[index].DrugCode)
		records[index].DrugName = strings.TrimSpace(records[index].DrugName)
		records[index].SettlementDate = normalizeDate(records[index].SettlementDate)
	}
}

func normalizeOCRText(value string) string {
	value = strings.ReplaceAll(value, "：", ":")
	value = strings.ReplaceAll(value, "，", " ")
	value = strings.ReplaceAll(value, ",", " ")
	value = strings.ReplaceAll(value, "年", "-")
	value = strings.ReplaceAll(value, "月", "-")
	value = strings.ReplaceAll(value, "日", "")
	value = strings.ReplaceAll(value, "_", "-")
	value = regexp.MustCompile(`\s+`).ReplaceAllString(value, " ")
	return strings.TrimSpace(value)
}

func matchFirst(value string, pattern string) string {
	matched := regexp.MustCompile(pattern).FindStringSubmatch(value)
	if len(matched) < 2 {
		return ""
	}
	return strings.TrimSpace(matched[1])
}

func normalizeDate(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = strings.ReplaceAll(value, "年", "-")
	value = strings.ReplaceAll(value, "月", "-")
	value = strings.ReplaceAll(value, "日", "")
	value = strings.ReplaceAll(value, "_", "-")
	datePattern := regexp.MustCompile(`\d{4}-\d{1,2}-\d{1,2}`)
	value = datePattern.FindString(value)
	if value == "" {
		return ""
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err == nil {
		return parsed.Format("2006-01-02")
	}
	parsed, err = time.Parse("2006-1-2", value)
	if err == nil {
		return parsed.Format("2006-01-02")
	}
	return ""
}

func truncateString(value string, maxLength int) string {
	if len(value) <= maxLength {
		return value
	}
	return value[:maxLength]
}

func getEnvOrDefault(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value != "" {
		return value
	}
	return strings.TrimSpace(fallback)
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

type validationError string

func newValidationError(message string) error {
	return validationError(message)
}

func (e validationError) Error() string {
	return string(e)
}
