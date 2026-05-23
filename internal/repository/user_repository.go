package repository

import (
	"context"
	"database/sql"

	"task-scheduler/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByToken(ctx context.Context, token string) (*model.User, error) {
	query := `SELECT id, username, password, display_name, token, status, created_at, updated_at FROM users WHERE token = ?`
	row := r.db.QueryRowContext(ctx, query, token)

	user := &model.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.DisplayName, &user.Token, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT id, username, password, display_name, token, status, created_at, updated_at FROM users WHERE username = ?`
	row := r.db.QueryRowContext(ctx, query, username)

	user := &model.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.DisplayName, &user.Token, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateToken(ctx context.Context, userID int64, token string) error {
	query := `UPDATE users SET token = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, token, userID)
	return err
}
