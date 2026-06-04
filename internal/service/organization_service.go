package service

import (
	"context"
	"fmt"

	"Tally-server/internal/model"
	"Tally-server/internal/repository"
)

type OrganizationService struct {
	repo repository.OrganizationRepository
}

func NewOrganizationService(repo repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

func (s *OrganizationService) Create(ctx context.Context, req model.CreateOrganizationRequest) (model.Organization, error) {
	if req.OrganizationName == "" {
		return model.Organization{}, fmt.Errorf("organization_name is required")
	}
	if req.Email == "" {
		return model.Organization{}, fmt.Errorf("email is required")
	}
	if req.Phone == "" {
		return model.Organization{}, fmt.Errorf("phone is required")
	}
	return s.repo.Create(ctx, req)
}

func (s *OrganizationService) GetByID(ctx context.Context, id string) (model.Organization, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrganizationService) List(ctx context.Context) ([]model.Organization, error) {
	return s.repo.List(ctx)
}

