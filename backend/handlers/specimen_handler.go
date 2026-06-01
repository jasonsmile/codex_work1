package handlers

import (
	"net/http"
	"strings"
	"time"

	"drug-info/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SpecimenHandler struct {
	db *gorm.DB
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

func inList(value string, options []string) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}
