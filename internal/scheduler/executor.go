package scheduler

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"task-scheduler/internal/model"
	"task-scheduler/internal/repository"
)

type Executor struct {
	execRepo *repository.ExecutionRepository
	taskRepo *repository.TaskRepository
}

func NewExecutor(execRepo *repository.ExecutionRepository, taskRepo *repository.TaskRepository) *Executor {
	return &Executor{
		execRepo: execRepo,
		taskRepo: taskRepo,
	}
}

func (e *Executor) Execute(ctx context.Context, task *model.Task, triggerType string) {
	now := time.Now()
	execLog := &model.ExecutionLog{
		TaskID:        task.ID,
		TaskName:      task.Name,
		TriggerType:   triggerType,
		Status:        "pending",
		StartTime:     &now,
		RetryAttempts: 0,
		Params:        buildParams(task),
	}

	logID, err := e.execRepo.Create(ctx, execLog)
	if err != nil {
		return
	}
	execLog.ID = logID

	execLog.Status = "running"
	execLog.StartTime = &now
	e.execRepo.Update(ctx, execLog)

	e.runWithRetry(ctx, task, execLog)
}

func (e *Executor) runWithRetry(ctx context.Context, task *model.Task, execLog *model.ExecutionLog) {
	maxAttempts := task.RetryCount + 1
	var lastErr error
	var output string

	for attempt := 0; attempt < maxAttempts; attempt++ {
		execLog.RetryAttempts = attempt

		output, lastErr = e.runCommand(ctx, task)

		if lastErr == nil {
			now := time.Now()
			execLog.Status = "success"
			execLog.EndTime = &now
			execLog.Duration = now.Sub(*execLog.StartTime).Milliseconds()
			execLog.Output = output
			execLog.ErrorMessage = ""
			e.execRepo.Update(ctx, execLog)
			return
		}

		if attempt < maxAttempts-1 {
			execLog.Status = "retrying"
			execLog.Output = output
			execLog.ErrorMessage = lastErr.Error()
			e.execRepo.Update(ctx, execLog)

			select {
			case <-time.After(time.Duration(task.RetryInterval) * time.Second):
			case <-ctx.Done():
				return
			}
		}
	}

	now := time.Now()
	execLog.EndTime = &now
	execLog.Duration = now.Sub(*execLog.StartTime).Milliseconds()
	execLog.Output = output
	if lastErr != nil {
		execLog.ErrorMessage = lastErr.Error()
	}
	execLog.Status = "failed"
	e.execRepo.Update(ctx, execLog)
}

func (e *Executor) runCommand(ctx context.Context, task *model.Task) (string, error) {
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

func buildParams(task *model.Task) string {
	return fmt.Sprintf("command=%s,timeout=%d,retry_count=%d,retry_interval=%d",
		task.Command, task.Timeout, task.RetryCount, task.RetryInterval)
}
