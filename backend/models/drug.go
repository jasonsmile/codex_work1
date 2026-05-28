package models

import "time"

type Drug struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name" gorm:"column:name;type:varchar(100);not null;index"`
	Manufacturer   string    `json:"manufacturer" gorm:"column:manufacturer;type:varchar(150)"`
	ApprovalNumber string    `json:"approvalNumber" gorm:"column:approval_number;type:varchar(100)"`
	Specification  string    `json:"specification" gorm:"column:specification;type:varchar(100)"`
	Price          float64   `json:"price" gorm:"column:price;type:decimal(10,2);not null"`
	Stock          int       `json:"stock" gorm:"column:stock;not null;check:stock > 0"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (Drug) TableName() string {
	return "drugs"
}

type CreateDrugRequest struct {
	Name           string  `json:"name" binding:"required"`
	Manufacturer   string  `json:"manufacturer"`
	ApprovalNumber string  `json:"approvalNumber"`
	Specification  string  `json:"specification"`
	Price          float64 `json:"price" binding:"required"`
	Stock          int     `json:"stock" binding:"required,min=1"`
}
