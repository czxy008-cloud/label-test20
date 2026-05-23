package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"task-scheduler/internal/model"
	"task-scheduler/internal/service"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req model.TaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "请求参数错误: "+err.Error()))
		return
	}

	task, err := h.taskService.Create(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "创建任务失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(task))
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "无效的任务ID"))
		return
	}

	var req model.TaskUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "请求参数错误: "+err.Error()))
		return
	}

	task, err := h.taskService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "更新任务失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(task))
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "无效的任务ID"))
		return
	}

	if err := h.taskService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "删除任务失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "无效的任务ID"))
		return
	}

	task, err := h.taskService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Error(404, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(task))
}

func (h *TaskHandler) List(c *gin.Context) {
	var req model.TaskListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "请求参数错误: "+err.Error()))
		return
	}

	resp, err := h.taskService.List(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "获取任务列表失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}
