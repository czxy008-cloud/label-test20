package repository

import (
	"context"
	"database/sql"
	"strings"

	"task-scheduler/internal/model"
)

type ExecutionRepository struct {
	db *sql.DB
}

func NewExecutionRepository(db *sql.DB) *ExecutionRepository {
	return &ExecutionRepository{db: db}
}

func (r *ExecutionRepository) Create(ctx context.Context, log *model.ExecutionLog) (int64, error) {
	query := `INSERT INTO execution_logs (task_id, task_name, trigger_type, status, start_time, retry_attempts, node_id) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		log.TaskID, log.TaskName, log.TriggerType, log.Status, log.StartTime, log.RetryAttempts, log.NodeID,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ExecutionRepository) Update(ctx context.Context, log *model.ExecutionLog) error {
	query := `UPDATE execution_logs SET status = ?, end_time = ?, duration = ?, output = ?, error_message = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		log.Status, log.EndTime, log.Duration, log.Output, log.ErrorMessage, log.ID,
	)
	return err
}

func (r *ExecutionRepository) GetByID(ctx context.Context, id int64) (*model.ExecutionLog, error) {
	query := `SELECT id, task_id, task_name, trigger_type, status, start_time, end_time, duration, retry_attempts, output, error_message, node_id, created_at FROM execution_logs WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	log := &model.ExecutionLog{}
	err := row.Scan(&log.ID, &log.TaskID, &log.TaskName, &log.TriggerType, &log.Status, &log.StartTime, &log.EndTime, &log.Duration, &log.RetryAttempts, &log.Output, &log.ErrorMessage, &log.NodeID, &log.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return log, nil
}

func (r *ExecutionRepository) List(ctx context.Context, req *model.ExecutionListRequest) ([]*model.ExecutionLog, int64, error) {
	var conditions []string
	var args []interface{}

	if req.TaskID > 0 {
		conditions = append(conditions, "task_id = ?")
		args = append(args, req.TaskID)
	}
	if req.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, req.Status)
	}
	if !req.StartTime.IsZero() {
		conditions = append(conditions, "start_time >= ?")
		args = append(args, req.StartTime)
	}
	if !req.EndTime.IsZero() {
		conditions = append(conditions, "start_time <= ?")
		args = append(args, req.EndTime)
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := `SELECT COUNT(*) FROM execution_logs ` + where
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, task_id, task_name, trigger_type, status, start_time, end_time, duration, retry_attempts, output, error_message, node_id, created_at FROM execution_logs ` + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, req.PageSize, (req.Page-1)*req.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*model.ExecutionLog
	for rows.Next() {
		log := &model.ExecutionLog{}
		err := rows.Scan(&log.ID, &log.TaskID, &log.TaskName, &log.TriggerType, &log.Status, &log.StartTime, &log.EndTime, &log.Duration, &log.RetryAttempts, &log.Output, &log.ErrorMessage, &log.NodeID, &log.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}
