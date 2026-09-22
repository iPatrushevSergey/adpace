package entity

import (
	"fmt"
	"time"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain"
)

const (
	CampaignStatusActive = "active"
	CampaignStatusPaused = "paused"

	PauseReasonBudgetDailyExceeded = "budget_daily_exceeded"
	PauseReasonBudgetTotalExceeded = "budget_total_exceeded"
	PauseReasonManual              = "manual"
)

type CampaignOption func(*Campaign)

type Campaign struct {
	CampaignID   string
	AdvertiserID string
	Name         string
	BudgetTotal  int64
	BudgetDaily  int64
	SpendTotal   int64
	SpendToday   int64
	SpendDay     time.Time
	Status       string
	PauseReason  *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewCampaign(
	campaignID, advertiserID, name string,
	budgetTotal, budgetDaily int64,
	updatedAt time.Time,
	opts ...CampaignOption,
) (Campaign, error) {
	if name == "" {
		return Campaign{}, fmt.Errorf("%w: name cannot be empty", domain.ErrBadInput)
	}
	if budgetTotal <= 0 {
		return Campaign{}, fmt.Errorf("%w: budget_total must be positive", domain.ErrBadInput)
	}
	if budgetDaily <= 0 {
		return Campaign{}, fmt.Errorf("%w: budget_daily must be positive", domain.ErrBadInput)
	}
	if budgetDaily > budgetTotal {
		return Campaign{}, fmt.Errorf("%w: budget_daily must not exceed budget_total", domain.ErrBadInput)
	}

	c := Campaign{
		CampaignID:   campaignID,
		AdvertiserID: advertiserID,
		Name:         name,
		BudgetTotal:  budgetTotal,
		BudgetDaily:  budgetDaily,
		Status:       CampaignStatusActive,
		UpdatedAt:    updatedAt,
	}

	for _, opt := range opts {
		opt(&c)
	}
	return c, nil
}

func WithCampaignCreatedAt(t time.Time) CampaignOption {
	return func(c *Campaign) { c.CreatedAt = t }
}

// Pause moves the campaign to paused state.
func (c *Campaign) Pause(reason string, updatedAt time.Time) error {
	if !IsValidPauseReason(reason) {
		return fmt.Errorf("%w: invalid pause reason %q", domain.ErrBadInput, reason)
	}
	if c.Status == CampaignStatusPaused {
		return fmt.Errorf("%w: campaign already paused", domain.ErrConflict)
	}
	c.Status = CampaignStatusPaused
	c.PauseReason = &reason
	c.UpdatedAt = updatedAt
	return nil
}

// Resume moves the campaign back to active state.
func (c *Campaign) Resume(updatedAt time.Time) error {
	if c.Status == CampaignStatusActive {
		return fmt.Errorf("%w: campaign already active", domain.ErrConflict)
	}
	c.Status = CampaignStatusActive
	c.PauseReason = nil
	c.UpdatedAt = updatedAt
	return nil
}

func IsValidPauseReason(reason string) bool {
	switch reason {
	case PauseReasonBudgetDailyExceeded, PauseReasonBudgetTotalExceeded, PauseReasonManual:
		return true
	}
	return false
}
