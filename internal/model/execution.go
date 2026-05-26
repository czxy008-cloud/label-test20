package model

import "time"

type ExecutionLog struct {
	ID            int64      `json:"id"`
	TaskID        int64      `json:"task_id"`
	TaskName      string     `json:"task_name"`
	TriggerType   string     `json:"trigger_type"`
	Status        string     `json:"status"`
	StartTime     *time.Time `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	Duration      int64      `json:"duration"`
	RetryAttempts int        `json:"retry_attempts"`
	Output        string     `json:"output"`
	ErrorMessage  string     `json:"error_message"`
	Params        string     `json:"params"`
	NodeID        string     `json:"node_id"`
	CreatedAt     time.Time  `json:"created_at"`
}

type ExecutionListRequest struct {
	TaskID    int64     `form:"task_id"`
	Status    string    `form:"status"`
	StartTime time.Time `form:"start_time"`
	EndTime   time.Time `form:"end_time"`
	Page      int       `form:"page,default=1"`
	PageSize  int       `form:"page_size,default=10"`
}

type ExecutionListResponse struct {
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	List     []*ExecutionLog `json:"list"`
}

type ExecuteTaskRequest struct {
	TaskID int64 `json:"task_id" binding:"required"`
}
