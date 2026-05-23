package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"task-scheduler/internal/model"
	"task-scheduler/internal/service"
)

type NodeHandler struct {
	nodeService *service.NodeService
}

func NewNodeHandler(nodeService *service.NodeService) *NodeHandler {
	return &NodeHandler{nodeService: nodeService}
}

func (h *NodeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	node, err := h.nodeService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Error(404, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(node))
}

func (h *NodeHandler) List(c *gin.Context) {
	var req model.NodeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "请求参数错误: "+err.Error()))
		return
	}

	resp, err := h.nodeService.List(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "获取节点列表失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}
