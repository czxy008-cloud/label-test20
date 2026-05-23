package service

import (
	"context"
	"errors"
	"time"

	"task-scheduler/internal/model"
	"task-scheduler/internal/repository"
)

type ExecutionService struct {
	execRepo *repository.ExecutionRepository
	taskRepo *repository.TaskRepository
}

func NewExecutionService(execRepo *repository.ExecutionRepository, taskRepo *repository.TaskRepository) *ExecutionService {
	return &ExecutionService{
		execRepo: execRepo,
		taskRepo: taskRepo,
	}
}

func (s *ExecutionService) TriggerTask(ctx context.Context, taskID int64) (*model.ExecutionLog, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("任务不存在")
	}
	if task.Status != 1 {
		return nil, errors.New("任务已禁用")
	}

	now := time.Now()
	log := &model.ExecutionLog{
		TaskID:        task.ID,
		TaskName:      task.Name,
		TriggerType:   "manual",
		Status:        "pending",
		StartTime:     &now,
		RetryAttempts: 0,
	}

	id, err := s.execRepo.Create(ctx, log)
	if err != nil {
		return nil, err
	}

	log.ID = id
	return log, nil
}

func (s *ExecutionService) UpdateExecution(ctx context.Context, log *model.ExecutionLog) error {
	return s.execRepo.Update(ctx, log)
}

func (s *ExecutionService) GetByID(ctx context.Context, id int64) (*model.ExecutionLog, error) {
	log, err := s.execRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if log == nil {
		return nil, errors.New("执行记录不存在")
	}
	return log, nil
}

func (s *ExecutionService) List(ctx context.Context, req *model.ExecutionListRequest) (*model.ExecutionListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	logs, total, err := s.execRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	return &model.ExecutionListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     logs,
	}, nil
}
