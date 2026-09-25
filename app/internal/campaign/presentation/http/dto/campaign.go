package dto

//go:generate go run github.com/mailru/easyjson/easyjson@v0.9.0 -all $GOFILE

import "time"

type CreateCampaignRequest struct {
	AdvertiserID string `json:"advertiser_id"`
	Name         string `json:"name"`
	BudgetTotal  int64  `json:"budget_total"`
	BudgetDaily  int64  `json:"budget_daily"`
}

type CreateCampaignResponse struct {
	ID           string    `json:"id"`
	AdvertiserID string    `json:"advertiser_id"`
	Name         string    `json:"name"`
	BudgetTotal  int64     `json:"budget_total"`
	BudgetDaily  int64     `json:"budget_daily"`
	SpendTotal   int64     `json:"spend_total"`
	SpendToday   int64     `json:"spend_today"`
	Status       string    `json:"status"`
	PauseReason  *string   `json:"pause_reason,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type GetCampaignResponse struct {
	ID           string    `json:"id"`
	AdvertiserID string    `json:"advertiser_id"`
	Name         string    `json:"name"`
	BudgetTotal  int64     `json:"budget_total"`
	BudgetDaily  int64     `json:"budget_daily"`
	SpendTotal   int64     `json:"spend_total"`
	SpendToday   int64     `json:"spend_today"`
	Status       string    `json:"status"`
	PauseReason  *string   `json:"pause_reason,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PatchCampaignRequest struct {
	Name        *string `json:"name,omitempty"`
	BudgetTotal *int64  `json:"budget_total,omitempty"`
	BudgetDaily *int64  `json:"budget_daily,omitempty"`
}

type PatchCampaignResponse struct {
	ID           string    `json:"id"`
	AdvertiserID string    `json:"advertiser_id"`
	Name         string    `json:"name"`
	BudgetTotal  int64     `json:"budget_total"`
	BudgetDaily  int64     `json:"budget_daily"`
	SpendTotal   int64     `json:"spend_total"`
	SpendToday   int64     `json:"spend_today"`
	Status       string    `json:"status"`
	PauseReason  *string   `json:"pause_reason,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PutCampaignRequest struct {
	Name        string `json:"name"`
	BudgetTotal int64  `json:"budget_total"`
	BudgetDaily int64  `json:"budget_daily"`
}

type PutCampaignResponse struct {
	ID           string    `json:"id"`
	AdvertiserID string    `json:"advertiser_id"`
	Name         string    `json:"name"`
	BudgetTotal  int64     `json:"budget_total"`
	BudgetDaily  int64     `json:"budget_daily"`
	SpendTotal   int64     `json:"spend_total"`
	SpendToday   int64     `json:"spend_today"`
	Status       string    `json:"status"`
	PauseReason  *string   `json:"pause_reason,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PauseCampaignRequest struct {
	Reason string `json:"reason"`
}

type PauseCampaignResponse struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	PauseReason *string   `json:"pause_reason,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ResumeCampaignResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}
