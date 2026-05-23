package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"task-scheduler/internal/model"
	"task-scheduler/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "请求参数错误: "+err.Error()))
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Error(401, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Error(401, "未登录"))
		return
	}

	c.JSON(http.StatusOK, model.Success(user))
}
