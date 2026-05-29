package handlers

import (
	"net/http"
	"strings"

	"drug-info/backend/models"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const defaultHeaderImg = "https://api.dicebear.com/10.x/bottts/png"

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(database *gorm.DB) *UserHandler {
	return &UserHandler{db: database}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateSysUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "请求参数不正确", err)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.HeaderImg = strings.TrimSpace(req.HeaderImg)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)

	if req.Username == "" {
		badRequest(c, "用户登录名不能为空", nil)
		return
	}
	if req.Password == "" {
		badRequest(c, "用户登录密码不能为空", nil)
		return
	}

	var existing models.SysUser
	err := h.db.Where("username = ?", req.Username).First(&existing).Error
	if err == nil {
		fail(c, http.StatusConflict, CodeConflict, "用户登录名已存在", nil)
		return
	}
	if err != gorm.ErrRecordNotFound {
		serverError(c, "检查用户登录名失败", err)
		return
	}

	userUUID := uuid.NewV4()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		serverError(c, "生成密码密文失败", err)
		return
	}

	if req.HeaderImg == "" {
		req.HeaderImg = defaultHeaderImg
	}
	if req.AuthorityID == 0 {
		req.AuthorityID = 888
	}
	if req.Enable == 0 {
		req.Enable = 1
	}

	user := models.SysUser{
		UUID:        userUUID.String(),
		Username:    req.Username,
		Password:    string(hashedPassword),
		HeaderImg:   req.HeaderImg,
		AuthorityID: req.AuthorityID,
		Phone:       req.Phone,
		Email:       req.Email,
		Enable:      req.Enable,
	}

	if err := h.db.Create(&user).Error; err != nil {
		serverError(c, "创建用户失败", err)
		return
	}

	success(c, http.StatusCreated, "创建用户成功", models.NewSysUserResponse(user))
}
