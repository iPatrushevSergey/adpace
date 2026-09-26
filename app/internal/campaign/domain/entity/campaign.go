package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
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
	CampaignID   uuid.UUID
	AdvertiserID uuid.UUID
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
	campaignID, advertiserID uuid.UUID,
	name string,
	budgetTotal, budgetDaily int64,
	updatedAt time.Time,
	opts ...CampaignOption,
) (Campaign, error) {
	c := Campaign{
		CampaignID:   campaignID,
		AdvertiserID: advertiserID,
		Status:       CampaignStatusActive,
		UpdatedAt:    updatedAt,
	}

	if err := c.SetName(&name); err != nil {
		return Campaign{}, err
	}
	if err := c.SetBudgets(&budgetTotal, &budgetDaily); err != nil {
		return Campaign{}, err
	}

	for _, opt := range opts {
		opt(&c)
	}
	return c, nil
}

func WithCampaignCreatedAt(t time.Time) CampaignOption {
	return func(c *Campaign) { c.CreatedAt = t }
}

func (c *Campaign) SetName(name *string) error {
	if name == nil {
		return nil
	}
	if !IsValidCampaignName(*name) {
		return fmt.Errorf("%w: name cannot be empty", domain.ErrBadInput)
	}
	c.Name = *name
	return nil
}

func (c *Campaign) SetBudgets(total, daily *int64) error {
	newTotal := c.BudgetTotal
	if total != nil {
		newTotal = *total
	}
	newDaily := c.BudgetDaily
	if daily != nil {
		newDaily = *daily
	}

	if !IsValidBudget(newTotal) {
		return fmt.Errorf("%w: budget_total must be positive", domain.ErrBadInput)
	}
	if !IsValidBudget(newDaily) {
		return fmt.Errorf("%w: budget_daily must be positive", domain.ErrBadInput)
	}
	if !IsValidBudgetPair(newDaily, newTotal) {
		return fmt.Errorf("%w: budget_daily must not exceed budget_total", domain.ErrBadInput)
	}

	c.BudgetTotal = newTotal
	c.BudgetDaily = newDaily
	return nil
}

func IsValidCampaignName(name string) bool {
	return name != ""
}

func IsValidBudget(amount int64) bool {
	return amount > 0
}

func IsValidBudgetPair(daily, total int64) bool {
	return daily <= total
}

func IsValidPauseReason(reason string) bool {
	switch reason {
	case PauseReasonBudgetDailyExceeded, PauseReasonBudgetTotalExceeded, PauseReasonManual:
		return true
	}
	return false
}
