package scheduler

import (
	"context"
	"log"
	"sync"

	"github.com/robfig/cron/v3"

	"task-scheduler/internal/model"
	"task-scheduler/internal/repository"
)

type Scheduler struct {
	cron       *cron.Cron
	taskRepo   *repository.TaskRepository
	execFunc   ExecutionFunc
	entries    map[int64]cron.EntryID
	mu         sync.RWMutex
}

type ExecutionFunc func(ctx context.Context, task *model.Task, triggerType string)

func NewScheduler(taskRepo *repository.TaskRepository, execFunc ExecutionFunc) *Scheduler {
	return &Scheduler{
		cron:     cron.New(cron.WithSeconds()),
		taskRepo: taskRepo,
		execFunc: execFunc,
		entries:  make(map[int64]cron.EntryID),
	}
}

func (s *Scheduler) Start() error {
	tasks, err := s.taskRepo.ListAllEnabled(context.Background())
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if err := s.registerTask(task); err != nil {
			log.Printf("调度器: 注册任务失败 taskID=%d name=%s err=%v", task.ID, task.Name, err)
		}
	}

	s.cron.Start()
	log.Printf("调度器: 已启动，注册任务数=%d", len(s.entries))
	return nil
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("调度器: 已停止")
}

func (s *Scheduler) registerTask(task *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task.CronExpression == "" {
		return nil
	}

	entryID, err := s.cron.AddFunc(task.CronExpression, func() {
		s.execFunc(context.Background(), task, "cron")
	})
	if err != nil {
		return err
	}

	s.entries[task.ID] = entryID
	log.Printf("调度器: 注册任务 taskID=%d name=%s cron=%s", task.ID, task.Name, task.CronExpression)
	return nil
}

func (s *Scheduler) ReloadTask(taskID int64) error {
	s.mu.Lock()
	if entryID, ok := s.entries[taskID]; ok {
		s.cron.Remove(entryID)
		delete(s.entries, taskID)
	}
	s.mu.Unlock()

	task, err := s.taskRepo.GetByID(context.Background(), taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return nil
	}

	if task.Status != 1 {
		log.Printf("调度器: 任务已禁用，跳过注册 taskID=%d", taskID)
		return nil
	}

	return s.registerTask(task)
}

func (s *Scheduler) RemoveTask(taskID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, ok := s.entries[taskID]; ok {
		s.cron.Remove(entryID)
		delete(s.entries, taskID)
		log.Printf("调度器: 移除任务 taskID=%d", taskID)
	}
}

func (s *Scheduler) TriggerTask(ctx context.Context, task *model.Task) {
	s.execFunc(ctx, task, "manual")
}
