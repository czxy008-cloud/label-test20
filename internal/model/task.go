package model

import "time"

type Task struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CronExpression string    `json:"cron_expression"`
	Command        string    `json:"command"`
	Timeout        int       `json:"timeout"`
	RetryCount     int       `json:"retry_count"`
	RetryInterval  int       `json:"retry_interval"`
	Status         int       `json:"status"`
	CreatedBy      int64     `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TaskCreateRequest struct {
	Name           string `json:"name" binding:"required,max=128"`
	Description    string `json:"description"`
	CronExpression string `json:"cron_expression" binding:"required"`
	Command        string `json:"command" binding:"required"`
	Timeout        int    `json:"timeout"`
	RetryCount     int    `json:"retry_count"`
	RetryInterval  int    `json:"retry_interval"`
	Status         int    `json:"status"`
}

type TaskUpdateRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	CronExpression string `json:"cron_expression"`
	Command        string `json:"command"`
	Timeout        int    `json:"timeout"`
	RetryCount     *int   `json:"retry_count"`
	RetryInterval  int    `json:"retry_interval"`
	Status         *int   `json:"status"`
}

type TaskListRequest struct {
	Name     string `form:"name"`
	Status   *int   `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type TaskListResponse struct {
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	List     []*Task `json:"list"`
}
