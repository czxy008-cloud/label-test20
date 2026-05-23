package repository

import (
	"context"
	"database/sql"
	"strings"

	"task-scheduler/internal/model"
)

type NodeRepository struct {
	db *sql.DB
}

func NewNodeRepository(db *sql.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

func (r *NodeRepository) Create(ctx context.Context, node *model.Node) error {
	query := `INSERT INTO nodes (id, name, ip_address, status, last_heartbeat, cpu_usage, memory_usage) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		node.ID, node.Name, node.IPAddress, node.Status, node.LastHeartbeat, node.CPUUsage, node.MemoryUsage,
	)
	return err
}

func (r *NodeRepository) UpdateStatus(ctx context.Context, nodeID string, status string) error {
	query := `UPDATE nodes SET status = ?, last_heartbeat = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, nodeID)
	return err
}

func (r *NodeRepository) GetByID(ctx context.Context, id string) (*model.Node, error) {
	query := `SELECT id, name, ip_address, status, last_heartbeat, cpu_usage, memory_usage, created_at, updated_at FROM nodes WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	node := &model.Node{}
	err := row.Scan(&node.ID, &node.Name, &node.IPAddress, &node.Status, &node.LastHeartbeat, &node.CPUUsage, &node.MemoryUsage, &node.CreatedAt, &node.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return node, nil
}

func (r *NodeRepository) List(ctx context.Context, req *model.NodeListRequest) ([]*model.Node, int64, error) {
	var conditions []string
	var args []interface{}

	if req.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, req.Status)
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := `SELECT COUNT(*) FROM nodes ` + where
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, ip_address, status, last_heartbeat, cpu_usage, memory_usage, created_at, updated_at FROM nodes ` + where + ` ORDER BY last_heartbeat DESC LIMIT ? OFFSET ?`
	args = append(args, req.PageSize, (req.Page-1)*req.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var nodes []*model.Node
	for rows.Next() {
		node := &model.Node{}
		err := rows.Scan(&node.ID, &node.Name, &node.IPAddress, &node.Status, &node.LastHeartbeat, &node.CPUUsage, &node.MemoryUsage, &node.CreatedAt, &node.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		nodes = append(nodes, node)
	}

	return nodes, total, nil
}

func (r *NodeRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM nodes WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
