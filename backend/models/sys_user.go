package models

import "time"

type SysUser struct {
	ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID        string     `json:"uuid" gorm:"column:uuid;type:varchar(191);index;comment:用户UUID"`
	Username    string     `json:"username" gorm:"column:username;type:varchar(191);index;comment:用户登录名"`
	Password    string     `json:"password" gorm:"column:password;type:varchar(191);comment:用户登录密码"`
	HeaderImg   string     `json:"headerImg" gorm:"column:header_img;type:varchar(191);default:https://api.dicebear.com/10.x/bottts/png;comment:用户头像"`
	AuthorityID uint64     `json:"authorityId" gorm:"column:authority_id;default:888;comment:用户角色ID"`
	Phone       string     `json:"phone" gorm:"column:phone;type:varchar(191);comment:用户手机号"`
	Email       string     `json:"email" gorm:"column:email;type:varchar(191);comment:用户邮箱"`
	Enable      int64      `json:"enable" gorm:"column:enable;default:1;comment:用户是否被冻结 1正常 2冻结"`
	CreatedAt   *time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt   *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

func (SysUser) TableName() string {
	return "sys_users"
}

type CreateSysUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	HeaderImg   string `json:"headerImg"`
	AuthorityID uint64 `json:"authorityId"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Enable      int64  `json:"enable"`
}

type SysUserResponse struct {
	ID          uint64     `json:"id"`
	UUID        string     `json:"uuid"`
	Username    string     `json:"username"`
	HeaderImg   string     `json:"headerImg"`
	AuthorityID uint64     `json:"authorityId"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	Enable      int64      `json:"enable"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

func NewSysUserResponse(user SysUser) SysUserResponse {
	return SysUserResponse{
		ID:          user.ID,
		UUID:        user.UUID,
		Username:    user.Username,
		HeaderImg:   user.HeaderImg,
		AuthorityID: user.AuthorityID,
		Phone:       user.Phone,
		Email:       user.Email,
		Enable:      user.Enable,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
