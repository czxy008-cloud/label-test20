package model

import "time"

type Node struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	IPAddress     string     `json:"ip_address"`
	Status        string     `json:"status"`
	LastHeartbeat *time.Time `json:"last_heartbeat"`
	CPUUsage      float64    `json:"cpu_usage"`
	MemoryUsage   float64    `json:"memory_usage"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type NodeListRequest struct {
	Status   string `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type NodeListResponse struct {
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	List     []*Node `json:"list"`
}
