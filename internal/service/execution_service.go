package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
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
	execLog := &model.ExecutionLog{
		TaskID:        task.ID,
		TaskName:      task.Name,
		TriggerType:   "manual",
		Status:        "pending",
		StartTime:     &now,
		RetryAttempts: 0,
		Params:        fmt.Sprintf("command=%s,timeout=%d,retry_count=%d,retry_interval=%d", task.Command, task.Timeout, task.RetryCount, task.RetryInterval),
	}

	id, err := s.execRepo.Create(ctx, execLog)
	if err != nil {
		return nil, err
	}
	execLog.ID = id

	go s.executeTask(task, execLog)

	return execLog, nil
}

func (s *ExecutionService) executeTask(task *model.Task, execLog *model.ExecutionLog) {
	ctx := context.Background()

	startTime := time.Now()
	execLog.Status = "running"
	execLog.StartTime = &startTime
	execLog.RetryAttempts = 0
	s.execRepo.Update(ctx, execLog)

	output, err := s.runCommand(ctx, task)
	if err == nil {
		now := time.Now()
		execLog.Status = "success"
		execLog.EndTime = &now
		execLog.Duration = now.Sub(startTime).Milliseconds()
		execLog.Output = output
		execLog.ErrorMessage = ""
		s.execRepo.Update(ctx, execLog)
		return
	}

	if task.RetryCount > 0 {
		execLog.Status = "retrying"
		execLog.Output = output
		execLog.ErrorMessage = err.Error()
		s.execRepo.Update(ctx, execLog)

		for attempt := 1; attempt <= task.RetryCount; attempt++ {
			execLog.RetryAttempts = attempt

			select {
			case <-time.After(time.Duration(task.RetryInterval) * time.Second):
			case <-ctx.Done():
				return
			}

			retryOutput, retryErr := s.runCommand(ctx, task)
			execLog.Output = retryOutput

			if retryErr == nil {
				now := time.Now()
				execLog.Status = "success"
				execLog.EndTime = &now
				execLog.Duration = now.Sub(startTime).Milliseconds()
				execLog.ErrorMessage = ""
				s.execRepo.Update(ctx, execLog)
				return
			}

			execLog.ErrorMessage = retryErr.Error()
			if attempt < task.RetryCount {
				execLog.Status = "retrying"
				s.execRepo.Update(ctx, execLog)
			}
		}
	}

	if task.RetryCount == 0 {
		execLog.Output = output
		execLog.ErrorMessage = err.Error()
	}
	now := time.Now()
	execLog.Status = "failed"
	execLog.EndTime = &now
	execLog.Duration = now.Sub(startTime).Milliseconds()
	s.execRepo.Update(ctx, execLog)
}

func (s *ExecutionService) runCommand(ctx context.Context, task *model.Task) (string, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(task.Timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(timeoutCtx, "cmd", "/C", task.Command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n[stderr]\n" + stderr.String()
	}

	if err != nil {
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return output, fmt.Errorf("任务执行超时(超过%d秒): %w", task.Timeout, err)
		}
		return output, fmt.Errorf("任务执行失败: %w", err)
	}

	return output, nil
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
