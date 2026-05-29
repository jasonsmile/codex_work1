package handlers

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func RBACMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawToken := strings.TrimSpace(c.GetHeader("Authorization"))
		rawToken = strings.TrimPrefix(rawToken, "Bearer ")
		if rawToken == "" {
			fail(c, http.StatusUnauthorized, 401, "请先登录", nil)
			c.Abort()
			return
		}

		claims, err := parseToken(rawToken)
		if err != nil {
			fail(c, http.StatusUnauthorized, 401, "登录状态无效", err)
			c.Abort()
			return
		}

		allowed, err := enforcer.Enforce(roleName(claims.AuthorityID), c.FullPath(), c.Request.Method)
		if err != nil {
			serverError(c, "权限校验失败", err)
			c.Abort()
			return
		}
		if !allowed {
			fail(c, http.StatusForbidden, 403, "无权访问该功能", nil)
			c.Abort()
			return
		}

		c.Set("userUUID", claims.UUID)
		c.Set("username", claims.Username)
		c.Set("authorityID", claims.AuthorityID)
		c.Next()
	}
}
