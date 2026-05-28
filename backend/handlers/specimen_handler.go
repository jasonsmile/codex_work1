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
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "请求参数不正确",
			"error":   err.Error(),
		})
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "性别必须为男或女"})
		return
	}
	if !inList(req.SampleType, []string{"组织", "血浆"}) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "送检标本类型不正确"})
		return
	}
	if !inList(req.PathologyType, []string{"腺癌", "鳞癌", "腺鳞癌", "大细胞神经内分泌癌", "小细胞肺癌", "其他"}) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "病理类型不正确"})
		return
	}
	if req.PDL1Expression < 0 || req.PDL1Expression > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "PD-L1表达必须在 0 到 100 之间"})
		return
	}
	if !inList(req.Stage, []string{"I", "II", "III", "IV"}) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "分期不正确"})
		return
	}
	if _, err := time.Parse("2006-01-02", req.InspectionDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "送检日期格式必须为 YYYY-MM-DD"})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "保存申请单失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "保存成功",
		"data":    application,
	})
}

func (h *SpecimenHandler) ListApplications(c *gin.Context) {
	var applications []models.SpecimenApplication
	if err := h.db.Order("created_at DESC").Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询申请单失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "查询成功",
		"data":    applications,
	})
}

func inList(value string, options []string) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}
