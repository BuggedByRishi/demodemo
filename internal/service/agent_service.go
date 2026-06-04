package service

import (
	"context"
	"crypto/sha256"
	"fmt"

	"Tally-server/internal/model"
	"Tally-server/internal/repository"
)

type AgentService struct {
	repo repository.AgentRepository
}

func NewAgentService(repo repository.AgentRepository) *AgentService {
	return &AgentService{repo: repo}
}

func (s *AgentService) Create(ctx context.Context, req model.CreateAgentRequest) (model.Agent, error) {
	if req.AgentName == "" {
		return model.Agent{}, fmt.Errorf("agent_name is required")
	}
	if req.Email == "" {
		return model.Agent{}, fmt.Errorf("email is required")
	}
	if req.GuacamolePassword == "" {
		return model.Agent{}, fmt.Errorf("guacamole_password is required")
	}

	// Hash the guacamole password before storing (SHA-256, matches Guacamole's default)
	h := sha256.Sum256([]byte(req.GuacamolePassword))

	return s.repo.Create(ctx, req, h[:])
}

func (s *AgentService) GetByID(ctx context.Context, id string) (model.Agent, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AgentService) ListByOrganization(ctx context.Context, orgID string) ([]model.Agent, error) {
	return s.repo.ListByOrganization(ctx, orgID)
}

