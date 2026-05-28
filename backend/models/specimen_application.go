package models

import "time"

type SpecimenApplication struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	Name               string    `json:"name" gorm:"column:name;type:varchar(50);not null;index"`
	Gender             string    `json:"gender" gorm:"column:gender;type:varchar(10);not null"`
	Age                int       `json:"age" gorm:"column:age;not null;check:age > 0"`
	IDNumber           string    `json:"idNumber" gorm:"column:id_number;type:varchar(40);not null;index"`
	SampleType         string    `json:"sampleType" gorm:"column:sample_type;type:varchar(20);not null"`
	PathologyType      string    `json:"pathologyType" gorm:"column:pathology_type;type:varchar(50);not null"`
	PDL1Expression     int       `json:"pdl1Expression" gorm:"column:pdl1_expression;not null"`
	DriverGeneMutation string    `json:"driverGeneMutation" gorm:"column:driver_gene_mutation;type:varchar(255)"`
	Stage              string    `json:"stage" gorm:"column:stage;type:varchar(10);not null"`
	LastTreatment      string    `json:"lastTreatment" gorm:"column:last_treatment;type:varchar(255)"`
	FollowUpTreatment  string    `json:"followUpTreatment" gorm:"column:follow_up_treatment;type:varchar(255)"`
	Doctor             string    `json:"doctor" gorm:"column:doctor;type:varchar(50);not null"`
	InspectionDate     string    `json:"inspectionDate" gorm:"column:inspection_date;type:varchar(10);not null"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

func (SpecimenApplication) TableName() string {
	return "specimen_applications"
}

type CreateSpecimenApplicationRequest struct {
	Name               string `json:"name" binding:"required"`
	Gender             string `json:"gender" binding:"required"`
	Age                int    `json:"age" binding:"required,min=1"`
	IDNumber           string `json:"idNumber" binding:"required"`
	SampleType         string `json:"sampleType" binding:"required"`
	PathologyType      string `json:"pathologyType" binding:"required"`
	PDL1Expression     int    `json:"pdl1Expression"`
	DriverGeneMutation string `json:"driverGeneMutation"`
	Stage              string `json:"stage" binding:"required"`
	LastTreatment      string `json:"lastTreatment"`
	FollowUpTreatment  string `json:"followUpTreatment"`
	Doctor             string `json:"doctor" binding:"required"`
	InspectionDate     string `json:"inspectionDate" binding:"required"`
}
