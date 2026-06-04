package repository

import (
	"context"
	"errors"
	"fmt"

	"Tally-server/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AgentRepository interface {
	Create(ctx context.Context, req model.CreateAgentRequest, hashedPassword []byte) (model.Agent, error)
	GetByID(ctx context.Context, id string) (model.Agent, error)
	ListByOrganization(ctx context.Context, orgID string) ([]model.Agent, error)
	SetActive(ctx context.Context, id string, active bool) error
}

type agentRepo struct{ db *pgxpool.Pool }

func NewAgentRepository(db *pgxpool.Pool) AgentRepository {
	return &agentRepo{db: db}
}

func (r *agentRepo) Create(ctx context.Context, req model.CreateAgentRequest, hashedPassword []byte) (model.Agent, error) {
	query := `
		INSERT INTO core.agents (organization_id, agent_name, phone, email, guacamole_password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING agent_id, organization_id, agent_name, phone, email,
		          guacamole_password, is_active, created_at, updated_at`

	var a model.Agent
	err := r.db.QueryRow(ctx, query,
		req.OrganizationID, req.AgentName, req.Phone, req.Email, hashedPassword,
	).Scan(
		&a.AgentID, &a.OrganizationID, &a.AgentName,
		&a.Phone, &a.Email, &a.GuacamolePassword,
		&a.IsActive, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return model.Agent{}, fmt.Errorf("create agent: %w", err)
	}
	return a, nil
}

func (r *agentRepo) GetByID(ctx context.Context, id string) (model.Agent, error) {
	query := `
		SELECT agent_id, organization_id, agent_name, phone, email,
		       guacamole_password, is_active, created_at, updated_at
		FROM core.agents WHERE agent_id = $1`

	var a model.Agent
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.AgentID, &a.OrganizationID, &a.AgentName,
		&a.Phone, &a.Email, &a.GuacamolePassword,
		&a.IsActive, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Agent{}, model.ErrNotFound
	}
	if err != nil {
		return model.Agent{}, fmt.Errorf("get agent: %w", err)
	}
	return a, nil
}

func (r *agentRepo) ListByOrganization(ctx context.Context, orgID string) ([]model.Agent, error) {
	query := `
		SELECT agent_id, organization_id, agent_name, phone, email,
		       guacamole_password, is_active, created_at, updated_at
		FROM core.agents WHERE organization_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("list agents: %w", err)
	}
	defer rows.Close()

	var agents []model.Agent
	for rows.Next() {
		var a model.Agent
		if err := rows.Scan(
			&a.AgentID, &a.OrganizationID, &a.AgentName,
			&a.Phone, &a.Email, &a.GuacamolePassword,
			&a.IsActive, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan agent: %w", err)
		}
		agents = append(agents, a)
	}
	return agents, rows.Err()
}

func (r *agentRepo) SetActive(ctx context.Context, id string, active bool) error {
	query := `UPDATE core.agents SET is_active = $1, updated_at = NOW() WHERE agent_id = $2`
	res, err := r.db.Exec(ctx, query, active, id)
	if err != nil {
		return fmt.Errorf("set agent active: %w", err)
	}
	if res.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
