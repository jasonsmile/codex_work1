package handlers

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"drug-info/backend/models"

	"github.com/extrame/xls"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

const maxSpecimenUploadSize = 10 << 20

type SpecimenHandler struct {
	db *gorm.DB
}

type importRowError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

type importSpecimenResult struct {
	TotalRows    int              `json:"totalRows"`
	SuccessCount int              `json:"successCount"`
	SkippedCount int              `json:"skippedCount"`
	Errors       []importRowError `json:"errors"`
}

func NewSpecimenHandler(database *gorm.DB) *SpecimenHandler {
	return &SpecimenHandler{db: database}
}

func (h *SpecimenHandler) CreateApplication(c *gin.Context) {
	var req models.CreateSpecimenApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数不正确", err)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Gender = strings.TrimSpace(req.Gender)
	req.IDNumber = strings.TrimSpace(req.IDNumber)
	req.SampleType = strings.TrimSpace(req.SampleType)
	req.PathologyType = strings.TrimSpace(req.PathologyType)
	req.DriverGeneMutation = strings.TrimSpace(req.DriverGeneMutation)
	req.Stage = strings.TrimSpace(req.Stage)
	req.LastTreatment = strings.TrimSpace(req.LastTreatment)
	req.FollowUpTreatment = strings.TrimSpace(req.FollowUpTreatment)
	req.Doctor = strings.TrimSpace(req.Doctor)
	req.InspectionDate = strings.TrimSpace(req.InspectionDate)

	if !inList(req.Gender, []string{"男", "女"}) {
		badRequest(c, "性别必须为男或女", nil)
		return
	}
	if !inList(req.SampleType, []string{"组织", "血浆"}) {
		badRequest(c, "送检标本类型不正确", nil)
		return
	}
	if !inList(req.PathologyType, []string{"腺癌", "鳞癌", "腺鳞癌", "大细胞神经内分泌癌", "小细胞肺癌", "其他"}) {
		badRequest(c, "病理类型不正确", nil)
		return
	}
	if req.PDL1Expression < 0 || req.PDL1Expression > 100 {
		badRequest(c, "PD-L1表达必须在 0 到 100 之间", nil)
		return
	}
	if !inList(req.Stage, []string{"I", "II", "III", "IV"}) {
		badRequest(c, "分期不正确", nil)
		return
	}
	if _, err := time.Parse("2006-01-02", req.InspectionDate); err != nil {
		badRequest(c, "送检日期格式必须为 YYYY-MM-DD", nil)
		return
	}

	application := models.SpecimenApplication{
		Name:               req.Name,
		Gender:             req.Gender,
		Age:                req.Age,
		IDNumber:           req.IDNumber,
		SampleType:         req.SampleType,
		PathologyType:      req.PathologyType,
		PDL1Expression:     req.PDL1Expression,
		DriverGeneMutation: req.DriverGeneMutation,
		Stage:              req.Stage,
		LastTreatment:      req.LastTreatment,
		FollowUpTreatment:  req.FollowUpTreatment,
		Doctor:             req.Doctor,
		InspectionDate:     req.InspectionDate,
	}

	if err := h.db.Create(&application).Error; err != nil {
		serverError(c, "保存申请单失败", err)
		return
	}

	success(c, http.StatusCreated, "保存成功", application)
}

func (h *SpecimenHandler) ListApplications(c *gin.Context) {
	var applications []models.SpecimenApplication

	name := strings.TrimSpace(c.Query("name"))
	idNumber := strings.TrimSpace(c.Query("idNumber"))
	inspectionDateStart := strings.TrimSpace(c.Query("inspectionDateStart"))
	inspectionDateEnd := strings.TrimSpace(c.Query("inspectionDateEnd"))

	if inspectionDateStart != "" {
		if _, err := time.Parse("2006-01-02", inspectionDateStart); err != nil {
			badRequest(c, "送检开始日期格式必须为 YYYY-MM-DD", nil)
			return
		}
	}
	if inspectionDateEnd != "" {
		if _, err := time.Parse("2006-01-02", inspectionDateEnd); err != nil {
			badRequest(c, "送检结束日期格式必须为 YYYY-MM-DD", nil)
			return
		}
	}
	if inspectionDateStart != "" && inspectionDateEnd != "" && inspectionDateStart > inspectionDateEnd {
		badRequest(c, "送检开始日期不能晚于结束日期", nil)
		return
	}

	query := h.db.Model(&models.SpecimenApplication{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if idNumber != "" {
		query = query.Where("id_number LIKE ?", "%"+idNumber+"%")
	}
	if inspectionDateStart != "" {
		query = query.Where("inspection_date >= ?", inspectionDateStart)
	}
	if inspectionDateEnd != "" {
		query = query.Where("inspection_date <= ?", inspectionDateEnd)
	}

	if err := query.Order("created_at DESC").Find(&applications).Error; err != nil {
		serverError(c, "查询申请单失败", err)
		return
	}

	success(c, http.StatusOK, "查询成功", applications)
}

func (h *SpecimenHandler) PreviewImportApplications(c *gin.Context) {
	result, _, ok := parseSpecimenImportFile(c)
	if !ok {
		return
	}

	success(c, http.StatusOK, "预览成功", result)
}

func (h *SpecimenHandler) ImportApplications(c *gin.Context) {
	result, applications, ok := parseSpecimenImportFile(c)
	if !ok {
		return
	}

	if len(applications) > 0 {
		if err := h.db.Create(&applications).Error; err != nil {
			serverError(c, "批量保存申请单失败", err)
			return
		}
	}

	result.SuccessCount = len(applications)
	success(c, http.StatusOK, "导入完成", result)
}

func parseSpecimenImportFile(c *gin.Context) (importSpecimenResult, []models.SpecimenApplication, bool) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSpecimenUploadSize+1024)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		badRequest(c, "请选择要上传的 Excel 文件", err)
		return importSpecimenResult{}, nil, false
	}
	if fileHeader.Size > maxSpecimenUploadSize {
		badRequest(c, "文件大小不能超过 10MB", nil)
		return importSpecimenResult{}, nil, false
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		badRequest(c, "文件格式仅支持 .xlsx、.xls", nil)
		return importSpecimenResult{}, nil, false
	}

	file, err := fileHeader.Open()
	if err != nil {
		serverError(c, "读取上传文件失败", err)
		return importSpecimenResult{}, nil, false
	}
	defer file.Close()

	var rows [][]string
	switch ext {
	case ".xlsx":
		rows, err = readXLSXRows(file)
	case ".xls":
		rows, err = readXLSRows(file)
	}
	if err != nil {
		badRequest(c, "解析 Excel 文件失败", err)
		return importSpecimenResult{}, nil, false
	}
	if len(rows) < 2 {
		badRequest(c, "Excel 文件至少需要包含表头和一行数据", nil)
		return importSpecimenResult{}, nil, false
	}

	headerMap := buildHeaderMap(rows[0])
	result := importSpecimenResult{Errors: make([]importRowError, 0)}
	applications := make([]models.SpecimenApplication, 0, len(rows)-1)

	for i, row := range rows[1:] {
		if isEmptyRow(row) {
			continue
		}
		result.TotalRows++
		application, err := parseSpecimenImportRow(row, headerMap)
		if err != nil {
			result.SkippedCount++
			result.Errors = append(result.Errors, importRowError{Row: i + 2, Message: err.Error()})
			continue
		}
		applications = append(applications, application)
	}

	result.SuccessCount = len(applications)
	return result, applications, true
}

func readXLSXRows(file io.Reader) ([][]string, error) {
	workbook, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer workbook.Close()

	sheets := workbook.GetSheetList()
	if len(sheets) == 0 {
		return nil, nil
	}
	return workbook.GetRows(sheets[0])
}

func readXLSRows(file io.Reader) ([][]string, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	workbook, err := xls.OpenReader(bytes.NewReader(data), "utf-8")
	if err != nil {
		return nil, err
	}
	return workbook.ReadAllCells(100000), nil
}

func buildHeaderMap(headers []string) map[string]int {
	headerMap := make(map[string]int)
	for index, header := range headers {
		key := normalizeHeader(header)
		if key != "" {
			headerMap[key] = index
		}
	}
	return headerMap
}

func normalizeHeader(header string) string {
	header = strings.TrimSpace(header)
	header = strings.ReplaceAll(header, " ", "")
	header = strings.ReplaceAll(header, "_", "")
	header = strings.ReplaceAll(header, "-", "")
	header = strings.ToLower(header)

	aliases := map[string]string{
		"姓名":                 "name",
		"name":               "name",
		"性别":                 "gender",
		"gender":             "gender",
		"年龄":                 "age",
		"age":                "age",
		"id号":                "id_number",
		"id号码":               "id_number",
		"idnumber":           "id_number",
		"送检标本类型":             "sample_type",
		"标本类型":               "sample_type",
		"sampletype":         "sample_type",
		"病理类型":               "pathology_type",
		"pathologytype":      "pathology_type",
		"pdl1表达":             "pdl1_expression",
		"pd-l1表达":            "pdl1_expression",
		"pdl1expression":     "pdl1_expression",
		"驱动基因突变":             "driver_gene_mutation",
		"驱动基因":               "driver_gene_mutation",
		"drivergenemutation": "driver_gene_mutation",
		"分期":                 "stage",
		"stage":              "stage",
		"末次治疗":               "last_treatment",
		"lasttreatment":      "last_treatment",
		"后续治疗方案":             "follow_up_treatment",
		"followuptreatment":  "follow_up_treatment",
		"送检医师":               "doctor",
		"doctor":             "doctor",
		"送检日期":               "inspection_date",
		"inspectiondate":     "inspection_date",
	}

	if value, ok := aliases[header]; ok {
		return value
	}
	return header
}

func parseSpecimenImportRow(row []string, headerMap map[string]int) (models.SpecimenApplication, error) {
	requiredFields := []string{
		"name", "gender", "age", "id_number", "sample_type", "pathology_type",
		"stage", "doctor", "inspection_date",
	}
	for _, field := range requiredFields {
		if _, ok := headerMap[field]; !ok {
			return models.SpecimenApplication{}, simpleError("缺少必填表头: " + field)
		}
	}

	name := getImportCell(row, headerMap, "name")
	gender := getImportCell(row, headerMap, "gender")
	ageText := getImportCell(row, headerMap, "age")
	idNumber := getImportCell(row, headerMap, "id_number")
	sampleType := getImportCell(row, headerMap, "sample_type")
	pathologyType := getImportCell(row, headerMap, "pathology_type")
	pdl1Text := getImportCell(row, headerMap, "pdl1_expression")
	driverGeneMutation := getImportCell(row, headerMap, "driver_gene_mutation")
	stage := normalizeStage(getImportCell(row, headerMap, "stage"))
	lastTreatment := getImportCell(row, headerMap, "last_treatment")
	followUpTreatment := getImportCell(row, headerMap, "follow_up_treatment")
	doctor := getImportCell(row, headerMap, "doctor")
	inspectionDate := getImportCell(row, headerMap, "inspection_date")

	if name == "" || gender == "" || ageText == "" || idNumber == "" || sampleType == "" || pathologyType == "" || stage == "" || doctor == "" || inspectionDate == "" {
		return models.SpecimenApplication{}, simpleError("必填字段不能为空")
	}

	age, err := strconv.Atoi(ageText)
	if err != nil || age <= 0 {
		return models.SpecimenApplication{}, simpleError("年龄必须为正整数")
	}

	pdl1Expression := 0
	if pdl1Text != "" {
		pdl1Text = strings.TrimSuffix(pdl1Text, "%")
		pdl1Expression, err = strconv.Atoi(pdl1Text)
		if err != nil || pdl1Expression < 0 || pdl1Expression > 100 {
			return models.SpecimenApplication{}, simpleError("PD-L1表达必须为 0 到 100 的整数")
		}
	}

	if !inList(gender, []string{"男", "女"}) {
		return models.SpecimenApplication{}, simpleError("性别必须为男或女")
	}
	if !inList(sampleType, []string{"组织", "血浆"}) {
		return models.SpecimenApplication{}, simpleError("送检标本类型必须为组织或血浆")
	}
	if !inList(pathologyType, []string{"腺癌", "鳞癌", "腺鳞癌", "大细胞神经内分泌癌", "小细胞肺癌", "其他"}) {
		return models.SpecimenApplication{}, simpleError("病理类型不正确")
	}
	if !inList(stage, []string{"I", "II", "III", "IV"}) {
		return models.SpecimenApplication{}, simpleError("分期必须为 I、II、III、IV")
	}
	if _, err := time.Parse("2006-01-02", inspectionDate); err != nil {
		return models.SpecimenApplication{}, simpleError("送检日期格式必须为 YYYY-MM-DD")
	}

	return models.SpecimenApplication{
		Name:               name,
		Gender:             gender,
		Age:                age,
		IDNumber:           idNumber,
		SampleType:         sampleType,
		PathologyType:      pathologyType,
		PDL1Expression:     pdl1Expression,
		DriverGeneMutation: driverGeneMutation,
		Stage:              stage,
		LastTreatment:      lastTreatment,
		FollowUpTreatment:  followUpTreatment,
		Doctor:             doctor,
		InspectionDate:     inspectionDate,
	}, nil
}

func getImportCell(row []string, headerMap map[string]int, key string) string {
	index, ok := headerMap[key]
	if !ok || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func normalizeStage(stage string) string {
	stage = strings.TrimSpace(strings.ToUpper(stage))
	switch stage {
	case "1":
		return "I"
	case "2":
		return "II"
	case "3":
		return "III"
	case "4":
		return "IV"
	default:
		return stage
	}
}

func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

type simpleError string

func (e simpleError) Error() string {
	return string(e)
}

func inList(value string, options []string) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}
