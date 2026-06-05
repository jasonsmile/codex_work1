package db

import (
	"drug-info/backend/models"

	"github.com/casbin/casbin/v2"
	casbinmodel "github.com/casbin/casbin/v2/model"
	"gorm.io/gorm"
)

var defaultCasbinRules = []models.CasbinRule{
	{PType: "g", V0: "role_888", V1: "admin"},
	{PType: "g", V0: "role_777", V1: "specimen_manager"},
	{PType: "g", V0: "role_999", V1: "viewer"},

	{PType: "p", V0: "admin", V1: "/api/drugs/add", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/api/drugs/get", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/api/specimens/add", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/api/specimens/import/preview", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/api/specimens/import", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/api/specimens/get", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/api/fileUploadAndDownload/upload", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/api/fileUploadAndDownload/get", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/api/fileUploadAndDownload/download/:id", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/api/files/upload", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/api/files/get", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/api/files/download/:id", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/api/trace_codes/recognize", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/api/trace_codes/confirm", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/api/trace_codes/get", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/api/trace_codes/delete", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/api/users/add", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/api/users/get", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/api/users/delete", V2: "POST"},

	{PType: "p", V0: "specimen_manager", V1: "/api/drugs/get", V2: "GET"},
	{PType: "p", V0: "specimen_manager", V1: "/api/specimens/add", V2: "POST"},
	{PType: "p", V0: "specimen_manager", V1: "/api/specimens/import/preview", V2: "POST"},
	{PType: "p", V0: "specimen_manager", V1: "/api/specimens/import", V2: "POST"},
	{PType: "p", V0: "specimen_manager", V1: "/api/specimens/get", V2: "GET"},
	{PType: "p", V0: "specimen_manager", V1: "/api/fileUploadAndDownload/upload", V2: "POST"},
	{PType: "p", V0: "specimen_manager", V1: "/api/fileUploadAndDownload/get", V2: "GET"},
	{PType: "p", V0: "specimen_manager", V1: "/api/fileUploadAndDownload/download/:id", V2: "GET"},
	{PType: "p", V0: "specimen_manager", V1: "/api/files/upload", V2: "POST"},
	{PType: "p", V0: "specimen_manager", V1: "/api/files/get", V2: "GET"},
	{PType: "p", V0: "specimen_manager", V1: "/api/files/download/:id", V2: "GET"},
	{PType: "p", V0: "specimen_manager", V1: "/api/trace_codes/recognize", V2: "POST"},
	{PType: "p", V0: "specimen_manager", V1: "/api/trace_codes/confirm", V2: "POST"},
	{PType: "p", V0: "specimen_manager", V1: "/api/trace_codes/get", V2: "GET"},
	{PType: "p", V0: "specimen_manager", V1: "/api/trace_codes/delete", V2: "POST"},

	{PType: "p", V0: "viewer", V1: "/api/drugs/get", V2: "GET"},
	{PType: "p", V0: "viewer", V1: "/api/specimens/get", V2: "GET"},
	{PType: "p", V0: "viewer", V1: "/api/fileUploadAndDownload/upload", V2: "POST"},
	{PType: "p", V0: "viewer", V1: "/api/fileUploadAndDownload/get", V2: "GET"},
	{PType: "p", V0: "viewer", V1: "/api/fileUploadAndDownload/download/:id", V2: "GET"},
	{PType: "p", V0: "viewer", V1: "/api/files/upload", V2: "POST"},
	{PType: "p", V0: "viewer", V1: "/api/files/get", V2: "GET"},
	{PType: "p", V0: "viewer", V1: "/api/files/download/:id", V2: "GET"},
	{PType: "p", V0: "viewer", V1: "/api/trace_codes/recognize", V2: "POST"},
	{PType: "p", V0: "viewer", V1: "/api/trace_codes/confirm", V2: "POST"},
	{PType: "p", V0: "viewer", V1: "/api/trace_codes/get", V2: "GET"},
	{PType: "p", V0: "viewer", V1: "/api/trace_codes/delete", V2: "POST"},
}

func NewRBACEnforcer(database *gorm.DB) (*casbin.Enforcer, error) {
	if err := seedCasbinRules(database); err != nil {
		return nil, err
	}

	model, err := casbinmodel.NewModelFromFile("rbac_model.conf")
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(model)
	if err != nil {
		return nil, err
	}

	var rules []models.CasbinRule
	if err := database.Order("id ASC").Find(&rules).Error; err != nil {
		return nil, err
	}

	for _, rule := range rules {
		values := casbinRuleValues(rule)
		switch rule.PType {
		case "p":
			if _, err := enforcer.AddPolicy(values...); err != nil {
				return nil, err
			}
		case "g":
			if _, err := enforcer.AddGroupingPolicy(values...); err != nil {
				return nil, err
			}
		}
	}

	return enforcer, nil
}

func seedCasbinRules(database *gorm.DB) error {
	for _, rule := range defaultCasbinRules {
		var existing models.CasbinRule
		err := database.Where(
			"ptype = ? AND v0 = ? AND v1 = ? AND v2 = ? AND v3 = ? AND v4 = ? AND v5 = ?",
			rule.PType,
			rule.V0,
			rule.V1,
			rule.V2,
			rule.V3,
			rule.V4,
			rule.V5,
		).First(&existing).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		if err := database.Create(&rule).Error; err != nil {
			return err
		}
	}
	return nil
}

func casbinRuleValues(rule models.CasbinRule) []interface{} {
	rawValues := []string{rule.V0, rule.V1, rule.V2, rule.V3, rule.V4, rule.V5}
	values := make([]interface{}, 0, len(rawValues))
	for _, value := range rawValues {
		if value == "" {
			continue
		}
		values = append(values, value)
	}
	return values
}
