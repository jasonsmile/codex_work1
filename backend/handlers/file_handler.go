package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"drug-info/backend/config"
	"drug-info/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"gorm.io/gorm"
)

const maxFileUploadSize = 5 << 20

var allowedUploadTypes = map[string]struct {
	FileType    string
	ContentType string
}{
	".jpg":  {FileType: "image", ContentType: "image/jpeg"},
	".jpeg": {FileType: "image", ContentType: "image/jpeg"},
	".png":  {FileType: "image", ContentType: "image/png"},
	".svg":  {FileType: "image", ContentType: "image/svg+xml"},
	".mp4":  {FileType: "video", ContentType: "video/mp4"},
	".txt":  {FileType: "text", ContentType: "text/plain"},
	".sql":  {FileType: "text", ContentType: "text/plain"},
	".xls":  {FileType: "text", ContentType: "application/vnd.ms-excel"},
	".xlsx": {FileType: "text", ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
	".doc":  {FileType: "text", ContentType: "application/msword"},
	".docx": {FileType: "text", ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
}

type FileHandler struct {
	db    *gorm.DB
	qiniu config.QiniuKodoConfig
}

type qiniuPutRet struct {
	Key  string `json:"key"`
	Hash string `json:"hash"`
}

func NewFileHandler(database *gorm.DB, qiniuConfig config.QiniuKodoConfig) *FileHandler {
	return &FileHandler{db: database, qiniu: qiniuConfig}
}

func (h *FileHandler) Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxFileUploadSize+1024*1024)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		badRequest(c, "请选择要上传的文件", err)
		return
	}
	if fileHeader.Size > maxFileUploadSize {
		badRequest(c, "文件大小不能超过 5MB", nil)
		return
	}

	fileName := strings.TrimSpace(fileHeader.Filename)
	ext := strings.ToLower(filepath.Ext(fileName))
	uploadType, ok := allowedUploadTypes[ext]
	if !ok {
		badRequest(c, "仅支持 jpg、png、svg 图片，mp4 视频，以及 txt、sql、xls、xlsx、doc、docx 文本文件", nil)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		serverError(c, "读取上传文件失败", err)
		return
	}
	defer file.Close()

	key := h.buildFileKey(fileName)
	putPolicy := storage.PutPolicy{Scope: h.qiniu.Bucket}
	mac := qbox.NewMac(h.qiniu.AccessKey, h.qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := qiniuPutRet{}
	putExtra := storage.PutExtra{}
	if err := formUploader.Put(context.Background(), &ret, upToken, key, file, fileHeader.Size, &putExtra); err != nil {
		serverError(c, "上传七牛云失败", err)
		return
	}
	if ret.Key == "" {
		ret.Key = key
	}

	displayName := strings.TrimSpace(c.PostForm("displayName"))
	if displayName == "" {
		displayName = strings.TrimSuffix(fileName, ext)
	}

	record := models.FileUploadAndDownload{
		DisplayName: displayName,
		FileName:    fileName,
		FileType:    uploadType.FileType,
		ContentType: uploadType.ContentType,
		FileSize:    fileHeader.Size,
		Key:         ret.Key,
		URL:         h.buildFileURL(ret.Key),
	}

	if err := h.db.Create(&record).Error; err != nil {
		serverError(c, "保存文件记录失败", err)
		return
	}

	success(c, http.StatusCreated, "上传成功", record)
}

func (h *FileHandler) List(c *gin.Context) {
	var files []models.FileUploadAndDownload
	if err := h.db.Order("created_at DESC").Find(&files).Error; err != nil {
		serverError(c, "查询文件列表失败", err)
		return
	}
	success(c, http.StatusOK, "查询成功", files)
}

func (h *FileHandler) Download(c *gin.Context) {
	var file models.FileUploadAndDownload
	if err := h.db.First(&file, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			badRequest(c, "文件不存在", nil)
			return
		}
		serverError(c, "查询文件失败", err)
		return
	}

	success(c, http.StatusOK, "查询成功", gin.H{
		"url":       h.buildPrivateDownloadURL(file),
		"key":       file.Key,
		"expiresAt": time.Now().Add(time.Hour).Format(time.RFC3339),
	})
}

func (h *FileHandler) buildFileKey(fileName string) string {
	cleanPath := strings.Trim(strings.TrimSpace(h.qiniu.Path), "/")
	name := sanitizeFileName(fileName)
	if name == "" {
		name = "upload"
	}
	key := fmt.Sprintf("%d_%s", time.Now().UnixNano(), name)
	if cleanPath == "" {
		return key
	}
	return path.Join(cleanPath, key)
}

func (h *FileHandler) buildFileURL(key string) string {
	domain := h.buildDomain()
	if domain == "" {
		return key
	}
	return domain + "/" + key
}

func (h *FileHandler) buildPrivateDownloadURL(file models.FileUploadAndDownload) string {
	query := url.Values{}
	if file.FileName != "" {
		query.Set("attname", file.FileName)
	}

	mac := qbox.NewMac(h.qiniu.AccessKey, h.qiniu.SecretKey)
	deadline := time.Now().Add(time.Hour).Unix()
	return storage.MakePrivateURLv2WithQuery(mac, h.buildDomain(), file.Key, query, deadline)
}

func (h *FileHandler) buildDomain() string {
	domain := strings.TrimRight(strings.TrimSpace(h.qiniu.Domain), "/")
	if domain == "" {
		return ""
	}
	if strings.HasPrefix(domain, "http://") || strings.HasPrefix(domain, "https://") {
		return domain
	}
	return "http://" + domain
}

func sanitizeFileName(fileName string) string {
	fileName = filepath.Base(fileName)
	fileName = strings.ReplaceAll(fileName, "\\", "_")
	fileName = strings.ReplaceAll(fileName, "/", "_")
	fileName = strings.ReplaceAll(fileName, " ", "_")
	return fileName
}
