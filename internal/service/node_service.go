package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"task-scheduler/internal/model"
	"task-scheduler/internal/repository"
)

type NodeService struct {
	nodeRepo *repository.NodeRepository
}

func NewNodeService(nodeRepo *repository.NodeRepository) *NodeService {
	return &NodeService{nodeRepo: nodeRepo}
}

func (s *NodeService) Register(ctx context.Context, node *model.Node) (*model.Node, error) {
	if node.ID == "" {
		node.ID = uuid.New().String()
	}
	if node.Status == "" {
		node.Status = "online"
	}

	existing, err := s.nodeRepo.GetByID(ctx, node.ID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("节点已存在")
	}

	if err := s.nodeRepo.Create(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

func (s *NodeService) Heartbeat(ctx context.Context, nodeID string, status string) error {
	node, err := s.nodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("节点不存在")
	}
	return s.nodeRepo.UpdateStatus(ctx, nodeID, status)
}

func (s *NodeService) GetByID(ctx context.Context, id string) (*model.Node, error) {
	node, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, errors.New("节点不存在")
	}
	return node, nil
}

func (s *NodeService) List(ctx context.Context, req *model.NodeListRequest) (*model.NodeListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	nodes, total, err := s.nodeRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	return &model.NodeListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     nodes,
	}, nil
}

func (s *NodeService) Delete(ctx context.Context, id string) error {
	node, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("节点不存在")
	}
	return s.nodeRepo.Delete(ctx, id)
}
