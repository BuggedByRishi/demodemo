package repository

import (
	"context"
	"errors"
	"fmt"

	"Tally-server/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReferenceRepository interface {
	ListCountries(ctx context.Context) ([]model.Country, error)
	ListStatesByCountry(ctx context.Context, countryID int) ([]model.State, error)
	ListThemes(ctx context.Context) ([]model.FrontendTheme, error)
	GetThemeByID(ctx context.Context, id string) (model.FrontendTheme, error)
}

type referenceRepo struct{ db *pgxpool.Pool }

func NewReferenceRepository(db *pgxpool.Pool) ReferenceRepository {
	return &referenceRepo{db: db}
}

func (r *referenceRepo) ListCountries(ctx context.Context) ([]model.Country, error) {
	rows, err := r.db.Query(ctx, `SELECT country_id, country_name FROM core.countries ORDER BY country_id`)
	if err != nil {
		return nil, fmt.Errorf("list countries: %w", err)
	}
	defer rows.Close()

	var countries []model.Country
	for rows.Next() {
		var c model.Country
		if err := rows.Scan(&c.CountryID, &c.CountryName); err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}
	return countries, rows.Err()
}

func (r *referenceRepo) ListStatesByCountry(ctx context.Context, countryID int) ([]model.State, error) {
	rows, err := r.db.Query(ctx,
		`SELECT state_id, country_id, state_name FROM core.states WHERE country_id = $1 ORDER BY state_id`,
		countryID,
	)
	if err != nil {
		return nil, fmt.Errorf("list states: %w", err)
	}
	defer rows.Close()

	var states []model.State
	for rows.Next() {
		var s model.State
		if err := rows.Scan(&s.StateID, &s.CountryID, &s.StateName); err != nil {
			return nil, err
		}
		states = append(states, s)
	}
	return states, rows.Err()
}

func (r *referenceRepo) ListThemes(ctx context.Context) ([]model.FrontendTheme, error) {
	rows, err := r.db.Query(ctx,
		`SELECT theme_id, theme_name, theme_object, is_active, created_at, updated_at FROM core.frontend_themes WHERE is_active = true ORDER BY theme_name`,
	)
	if err != nil {
		return nil, fmt.Errorf("list themes: %w", err)
	}
	defer rows.Close()

	var themes []model.FrontendTheme
	for rows.Next() {
		var t model.FrontendTheme
		if err := rows.Scan(&t.ThemeID, &t.ThemeName, &t.ThemeObject, &t.IsActive, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		themes = append(themes, t)
	}
	return themes, rows.Err()
}

func (r *referenceRepo) GetThemeByID(ctx context.Context, id string) (model.FrontendTheme, error) {
	var t model.FrontendTheme
	err := r.db.QueryRow(ctx,
		`SELECT theme_id, theme_name, theme_object, is_active, created_at, updated_at FROM core.frontend_themes WHERE theme_id = $1`,
		id,
	).Scan(&t.ThemeID, &t.ThemeName, &t.ThemeObject, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.FrontendTheme{}, model.ErrNotFound
	}
	if err != nil {
		return model.FrontendTheme{}, fmt.Errorf("get theme: %w", err)
	}
	return t, nil
}
