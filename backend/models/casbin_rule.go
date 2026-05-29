package models

type CasbinRule struct {
	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	PType string `json:"ptype" gorm:"column:ptype;type:varchar(16);not null;index"`
	V0    string `json:"v0" gorm:"column:v0;type:varchar(191);default:'';index"`
	V1    string `json:"v1" gorm:"column:v1;type:varchar(191);default:'';index"`
	V2    string `json:"v2" gorm:"column:v2;type:varchar(191);default:''"`
	V3    string `json:"v3" gorm:"column:v3;type:varchar(191);default:''"`
	V4    string `json:"v4" gorm:"column:v4;type:varchar(191);default:''"`
	V5    string `json:"v5" gorm:"column:v5;type:varchar(191);default:''"`
}

func (CasbinRule) TableName() string {
	return "casbin_rules"
}
