package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"task-scheduler/internal/model"
	"task-scheduler/internal/service"
)

type ExecutionHandler struct {
	execService *service.ExecutionService
}

func NewExecutionHandler(execService *service.ExecutionService) *ExecutionHandler {
	return &ExecutionHandler{execService: execService}
}

func (h *ExecutionHandler) TriggerTask(c *gin.Context) {
	var req model.ExecuteTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "请求参数错误: "+err.Error()))
		return
	}

	log, err := h.execService.TriggerTask(c.Request.Context(), req.TaskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "触发任务失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(log))
}

func (h *ExecutionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "无效的执行ID"))
		return
	}

	log, err := h.execService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Error(404, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(log))
}

func (h *ExecutionHandler) List(c *gin.Context) {
	var req model.ExecutionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "请求参数错误: "+err.Error()))
		return
	}

	resp, err := h.execService.List(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "获取执行列表失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}
