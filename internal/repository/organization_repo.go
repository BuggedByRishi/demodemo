package repository

import (
	"context"
	"errors"
	"fmt"

	"Tally-server/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationRepository interface {
	Create(ctx context.Context, req model.CreateOrganizationRequest) (model.Organization, error)
	GetByID(ctx context.Context, id string) (model.Organization, error)
	List(ctx context.Context) ([]model.Organization, error)
	SetActive(ctx context.Context, id string, active bool) error
}

type orgRepo struct{ db *pgxpool.Pool }

func NewOrganizationRepository(db *pgxpool.Pool) OrganizationRepository {
	return &orgRepo{db: db}
}

func (r *orgRepo) Create(ctx context.Context, req model.CreateOrganizationRequest) (model.Organization, error) {
	query := `
		INSERT INTO core.organizations
			(country_id, state_id, theme_id, phone, email, organization_name, city, address, logo_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING organization_id, country_id, state_id, theme_id, phone, email,
		          organization_name, city, address, logo_url, is_whitelabel, is_active, created_at, updated_at`

	var o model.Organization
	err := r.db.QueryRow(ctx, query,
		req.CountryID, req.StateID, req.ThemeID,
		req.Phone, req.Email, req.OrganizationName,
		req.City, req.Address, req.LogoURL,
	).Scan(
		&o.OrganizationID, &o.CountryID, &o.StateID, &o.ThemeID,
		&o.Phone, &o.Email, &o.OrganizationName,
		&o.City, &o.Address, &o.LogoURL,
		&o.IsWhitelabel, &o.IsActive, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return model.Organization{}, fmt.Errorf("create organization: %w", err)
	}
	return o, nil
}

func (r *orgRepo) GetByID(ctx context.Context, id string) (model.Organization, error) {
	query := `
		SELECT organization_id, country_id, state_id, theme_id, phone, email,
		       organization_name, city, address, logo_url, is_whitelabel, is_active, created_at, updated_at
		FROM core.organizations WHERE organization_id = $1`

	var o model.Organization
	err := r.db.QueryRow(ctx, query, id).Scan(
		&o.OrganizationID, &o.CountryID, &o.StateID, &o.ThemeID,
		&o.Phone, &o.Email, &o.OrganizationName,
		&o.City, &o.Address, &o.LogoURL,
		&o.IsWhitelabel, &o.IsActive, &o.CreatedAt, &o.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Organization{}, model.ErrNotFound
	}
	if err != nil {
		return model.Organization{}, fmt.Errorf("get organization: %w", err)
	}
	return o, nil
}

func (r *orgRepo) List(ctx context.Context) ([]model.Organization, error) {
	query := `
		SELECT organization_id, country_id, state_id, theme_id, phone, email,
		       organization_name, city, address, logo_url, is_whitelabel, is_active, created_at, updated_at
		FROM core.organizations ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}
	defer rows.Close()

	var orgs []model.Organization
	for rows.Next() {
		var o model.Organization
		if err := rows.Scan(
			&o.OrganizationID, &o.CountryID, &o.StateID, &o.ThemeID,
			&o.Phone, &o.Email, &o.OrganizationName,
			&o.City, &o.Address, &o.LogoURL,
			&o.IsWhitelabel, &o.IsActive, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan organization: %w", err)
		}
		orgs = append(orgs, o)
	}
	return orgs, rows.Err()
}

func (r *orgRepo) SetActive(ctx context.Context, id string, active bool) error {
	query := `UPDATE core.organizations SET is_active = $1, updated_at = NOW() WHERE organization_id = $2`
	res, err := r.db.Exec(ctx, query, active, id)
	if err != nil {
		return fmt.Errorf("set organization active: %w", err)
	}
	if res.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
