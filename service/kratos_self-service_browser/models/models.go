package models

import "time"

// ServiceFlowsRequest for parsed JSON data from http request
type ServiceFlowsRequest struct {
	ID         string    `json:"id"`
	ExpiresAt  time.Time `json:"expires_at"`
	IssuedAt   time.Time `json:"issued_at"`
	RequestURL string    `json:"request_url"`
	Methods    Methods   `json:"methods"`
}

// Methods struct
type Methods struct {
	OIDC     OIDC     `json:"oidc"`
	Password Password `json:"password"`
}

// OIDC struct
type OIDC struct {
	Method string `json:"method"`
	Config Config `json:"config"`
}

// Password struct
type Password struct {
	Method string `json:"method"`
	Config Config `json:"config"`
}

// Config struct
type Config struct {
	Action string   `json:"action"`
	Method string   `json:"method"`
	Fields []Fields `json:"fields"`
}

// Fields struct
type Fields struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Value    string `json:"value"`
}

// OryKratosRegistrationForm struct for kratos registration form field
type OryKratosRegistrationForm struct {
}
