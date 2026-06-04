package model

import (
	"encoding/json"
	"time"
)

// Country maps to core.countries
type Country struct {
	CountryID   int    `json:"country_id"`
	CountryName string `json:"country_name"`
}

// State maps to core.states
type State struct {
	StateID   int    `json:"state_id"`
	CountryID int    `json:"country_id"`
	StateName string `json:"state_name"`
}

// FrontendTheme maps to core.frontend_themes
type FrontendTheme struct {
	ThemeID     string          `json:"theme_id"`
	ThemeName   string          `json:"theme_name"`
	ThemeObject json.RawMessage `json:"theme_object"` // stored as JSONB
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty"`
}

// Partner maps to core.partners
type Partner struct {
	PartnerID   string     `json:"partner_id"`
	CountryID   int        `json:"country_id"`
	StateID     int        `json:"state_id"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	PartnerName string     `json:"partner_name"`
	City        string     `json:"city"`
	Address     string     `json:"address"`
	LogoURL     string     `json:"logo_url"`
	Domain      *string    `json:"domain,omitempty"`
	ThemeID     string     `json:"theme_id"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// Organization maps to core.organizations
type Organization struct {
	OrganizationID   string     `json:"organization_id"`
	CountryID        int        `json:"country_id"`
	StateID          int        `json:"state_id"`
	ThemeID          string     `json:"theme_id"`
	Phone            string     `json:"phone"`
	Email            string     `json:"email"`
	OrganizationName string     `json:"organization_name"`
	City             string     `json:"city"`
	Address          string     `json:"address"`
	LogoURL          *string    `json:"logo_url,omitempty"`
	IsWhitelabel     bool       `json:"is_whitelabel"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

// OrganizationTeam maps to core.organization_teams
type OrganizationTeam struct {
	TeamID         string     `json:"team_id"`
	OrganizationID string     `json:"organization_id"`
	Team           string     `json:"team"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

// Agent maps to core.agents
type Agent struct {
	AgentID            string     `json:"agent_id"`
	OrganizationID     string     `json:"organization_id"`
	AgentName          string     `json:"agent_name"`
	Phone              string     `json:"phone"`
	Email              string     `json:"email"`
	GuacamolePassword  []byte     `json:"-"` // BYTEA — never expose in JSON
	IsActive           bool       `json:"is_active"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
}

// --- Request types ---

type CreateOrganizationRequest struct {
	CountryID        int    `json:"country_id"`
	StateID          int    `json:"state_id"`
	ThemeID          string `json:"theme_id"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	OrganizationName string `json:"organization_name"`
	City             string `json:"city"`
	Address          string `json:"address"`
	LogoURL          string `json:"logo_url"`
}

type CreateAgentRequest struct {
	OrganizationID    string `json:"organization_id"`
	AgentName         string `json:"agent_name"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	GuacamolePassword string `json:"guacamole_password"` // plain text; hashed before insert
}

type CreatePartnerRequest struct {
	CountryID   int    `json:"country_id"`
	StateID     int    `json:"state_id"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	PartnerName string `json:"partner_name"`
	City        string `json:"city"`
	Address     string `json:"address"`
	LogoURL     string `json:"logo_url"`
	ThemeID     string `json:"theme_id"`
}

