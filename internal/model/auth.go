package model

import "time"

// User maps to auth.users
type User struct {
	ID                  string     `json:"id"`
	Name                string     `json:"name"`
	Image               *string    `json:"image,omitempty"`
	Email               *string    `json:"email,omitempty"`
	EmailVerified       *bool      `json:"email_verified,omitempty"`
	PhoneNumber         *string    `json:"phone_number,omitempty"`
	PhoneNumberVerified *bool      `json:"phone_number_verified,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
}

// Session maps to auth.sessions
type Session struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Token     string     `json:"token"`
	IPAddress *string    `json:"ip_address,omitempty"`
	UserAgent *string    `json:"user_agent,omitempty"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// Account maps to auth.accounts
type Account struct {
	ID                    string     `json:"id"`
	UserID                string     `json:"user_id"`
	AccountID             string     `json:"account_id"`
	ProviderID            string     `json:"provider_id"`
	AccessToken           *string    `json:"access_token,omitempty"`
	RefreshToken          *string    `json:"refresh_token,omitempty"`
	IDToken               *string    `json:"id_token,omitempty"`
	Scope                 *string    `json:"scope,omitempty"`
	AccessTokenExpiresAt  *time.Time `json:"access_token_expires_at,omitempty"`
	RefreshTokenExpiresAt *time.Time `json:"refresh_token_expires_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty"`
}

// Verification maps to auth.verifications
type Verification struct {
	ID         string     `json:"id"`
	Identifier string     `json:"identifier"`
	Value      string     `json:"value"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

