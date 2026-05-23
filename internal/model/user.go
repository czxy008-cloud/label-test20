package model

import "time"

type User struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"-"`
	DisplayName string    `json:"display_name"`
	Token       string    `json:"-"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
