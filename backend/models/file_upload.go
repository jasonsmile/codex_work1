package models

import "time"

type FileUploadAndDownload struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	DisplayName string    `json:"displayName" gorm:"column:display_name;type:varchar(255);not null"`
	FileName    string    `json:"fileName" gorm:"column:file_name;type:varchar(255);not null"`
	FileType    string    `json:"fileType" gorm:"column:file_type;type:varchar(20);not null;index"`
	ContentType string    `json:"contentType" gorm:"column:content_type;type:varchar(100);not null"`
	FileSize    int64     `json:"fileSize" gorm:"column:file_size;not null"`
	Key         string    `json:"key" gorm:"column:file_key;type:varchar(512);not null;uniqueIndex"`
	URL         string    `json:"url" gorm:"column:url;type:varchar(1024);not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (FileUploadAndDownload) TableName() string {
	return "file_upload_and_downloads"
}
