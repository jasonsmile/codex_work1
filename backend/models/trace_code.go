package models

import "time"

type TraceCode struct {
	ID                      uint      `json:"id" gorm:"primaryKey"`
	TransactionSerialNumber string    `json:"transaction_serial_number" gorm:"column:transaction_serial_number;type:varchar(120);not null;index"`
	DrugCode                string    `json:"drug_code" gorm:"column:drug_code;type:varchar(120);not null;index"`
	DrugName                string    `json:"drug_name" gorm:"column:drug_name;type:varchar(150);not null"`
	SettlementDate          string    `json:"settlement_date" gorm:"column:settlement_date;type:varchar(20);not null;index"`
	SourceFileName          string    `json:"source_file_name" gorm:"column:source_file_name;type:varchar(255)"`
	ContentType             string    `json:"content_type" gorm:"column:content_type;type:varchar(100)"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func (TraceCode) TableName() string {
	return "trace_codes"
}
