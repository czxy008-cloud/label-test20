package service

import (
	"context"
	"errors"

	"task-scheduler/internal/model"
	"task-scheduler/internal/repository"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
}

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) Create(ctx context.Context, req *model.TaskCreateRequest, userID int64) (*model.Task, error) {
	if req.Timeout <= 0 {
		req.Timeout = 300
	}
	if req.RetryCount < 0 {
		req.RetryCount = 0
	}
	if req.RetryInterval <= 0 {
		req.RetryInterval = 60
	}

	task := &model.Task{
		Name:           req.Name,
		Description:    req.Description,
		CronExpression: req.CronExpression,
		Command:        req.Command,
		Timeout:        req.Timeout,
		RetryCount:     req.RetryCount,
		RetryInterval:  req.RetryInterval,
		Status:         req.Status,
		CreatedBy:      userID,
	}

	id, err := s.taskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	task.ID = id
	return task, nil
}

func (s *TaskService) Update(ctx context.Context, id int64, req *model.TaskUpdateRequest) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("任务不存在")
	}

	if req.Name != "" {
		task.Name = req.Name
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.CronExpression != "" {
		task.CronExpression = req.CronExpression
	}
	if req.Command != "" {
		task.Command = req.Command
	}
	if req.Timeout > 0 {
		task.Timeout = req.Timeout
	}
	if req.RetryCount >= 0 {
		task.RetryCount = req.RetryCount
	}
	if req.RetryInterval > 0 {
		task.RetryInterval = req.RetryInterval
	}
	if req.Status != 0 {
		task.Status = req.Status
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("任务不存在")
	}
	return s.taskRepo.Delete(ctx, id)
}

func (s *TaskService) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("任务不存在")
	}
	return task, nil
}

func (s *TaskService) List(ctx context.Context, req *model.TaskListRequest) (*model.TaskListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	tasks, total, err := s.taskRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	return &model.TaskListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     tasks,
	}, nil
}
