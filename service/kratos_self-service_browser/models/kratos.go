package models

import "time"

// CurrentKratosSession models
type CurrentKratosSession struct {
	SID             string     `json:"sid"`
	AuthenticatedAt *time.Time `json:"authenticated_at"`
	ExpiresAt       *time.Time `json:"expires_at"`
	Identity        Identity   `json:"identity"`
	IssuedAt        *time.Time `json:"issued_at"`
}

// Identity models
type Identity struct {
	ID              string      `json:"id"`
	Traits          Traits      `json:"traits"`
	TraitsSchemaID  string      `json:"traits_schema_id"`
	TraitsSchemaURL string      `json:"traits_schema_url"`
	Addresses       []Addresses `json:"addresses"`
}

// Traits models
type Traits struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Addresses models
type Addresses struct {
	ExpiresAt  *time.Time `json:"expires_at"`
	ID         string     `json:"id"`
	Value      string     `json:"value"`
	Verified   bool       `json:"verified"`
	VerifiedAt *time.Time `json:"verified_at"`
	Via        string     `json:"via"`
}
