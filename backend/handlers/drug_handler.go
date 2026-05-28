package handlers

import (
	"net/http"
	"strings"

	"drug-info/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DrugHandler struct {
	db *gorm.DB
}

func NewDrugHandler(database *gorm.DB) *DrugHandler {
	return &DrugHandler{db: database}
}

func (h *DrugHandler) CreateDrug(c *gin.Context) {
	var req models.CreateDrugRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "请求参数不正确",
			"error":   err.Error(),
		})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Manufacturer = strings.TrimSpace(req.Manufacturer)
	req.ApprovalNumber = strings.TrimSpace(req.ApprovalNumber)
	req.Specification = strings.TrimSpace(req.Specification)

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "药品名称不能为空"})
		return
	}
	if req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "价格必须大于 0"})
		return
	}
	if req.Stock <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "库存数量必须大于 0"})
		return
	}

	drug := models.Drug{
		Name:           req.Name,
		Manufacturer:   req.Manufacturer,
		ApprovalNumber: req.ApprovalNumber,
		Specification:  req.Specification,
		Price:          req.Price,
		Stock:          req.Stock,
	}

	if err := h.db.Create(&drug).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "保存药品失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "保存成功",
		"data":    drug,
	})
}

func (h *DrugHandler) ListDrugs(c *gin.Context) {
	name := strings.TrimSpace(c.Query("name"))

	query := h.db.Model(&models.Drug{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	var drugs []models.Drug
	if err := query.Order("created_at DESC").Find(&drugs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询药品失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "查询成功",
		"data":    drugs,
	})
}
