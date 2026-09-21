package dto

//go:generate go run github.com/mailru/easyjson/easyjson@v0.9.0 -all $GOFILE

import "time"

type CreateAdvertiserRequest struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type CreateAdvertiserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetAdvertiserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PatchAdvertiserRequest struct {
	Name    *string `json:"name,omitempty"`
	Country *string `json:"country,omitempty"`
}

type PatchAdvertiserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PutAdvertiserRequest struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type PutAdvertiserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
