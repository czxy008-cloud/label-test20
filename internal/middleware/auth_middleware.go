package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"task-scheduler/internal/model"
	"task-scheduler/internal/service"
)

type AuthMiddleware struct {
	authService *service.AuthService
	headerName  string
}

func NewAuthMiddleware(authService *service.AuthService, headerName string) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		headerName:  headerName,
	}
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(m.headerName)
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, model.Error(401, "未提供认证令牌"))
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		user, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.Error(401, "认证失败: "+err.Error()))
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}
