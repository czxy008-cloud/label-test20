package repository

import (
	"context"
	"database/sql"
	"strings"

	"task-scheduler/internal/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) (int64, error) {
	query := `INSERT INTO tasks (name, description, cron_expression, command, timeout, retry_count, retry_interval, status, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		task.Name, task.Description, task.CronExpression, task.Command,
		task.Timeout, task.RetryCount, task.RetryInterval, task.Status, task.CreatedBy,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	query := `UPDATE tasks SET name = ?, description = ?, cron_expression = ?, command = ?, timeout = ?, retry_count = ?, retry_interval = ?, status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		task.Name, task.Description, task.CronExpression, task.Command,
		task.Timeout, task.RetryCount, task.RetryInterval, task.Status, task.ID,
	)
	return err
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	query := `SELECT id, name, description, cron_expression, command, timeout, retry_count, retry_interval, status, created_by, created_at, updated_at FROM tasks WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	task := &model.Task{}
	err := row.Scan(&task.ID, &task.Name, &task.Description, &task.CronExpression, &task.Command,
		&task.Timeout, &task.RetryCount, &task.RetryInterval, &task.Status, &task.CreatedBy, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return task, nil
}

func (r *TaskRepository) List(ctx context.Context, req *model.TaskListRequest) ([]*model.Task, int64, error) {
	var conditions []string
	var args []interface{}

	if req.Name != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+req.Name+"%")
	}
	if req.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *req.Status)
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := `SELECT COUNT(*) FROM tasks ` + where
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, description, cron_expression, command, timeout, retry_count, retry_interval, status, created_by, created_at, updated_at FROM tasks ` + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, req.PageSize, (req.Page-1)*req.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		task := &model.Task{}
		err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.CronExpression, &task.Command,
			&task.Timeout, &task.RetryCount, &task.RetryInterval, &task.Status, &task.CreatedBy, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}
