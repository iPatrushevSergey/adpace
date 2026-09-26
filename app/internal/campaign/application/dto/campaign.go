package dto

import "github.com/google/uuid"

type CreateCampaignInput struct {
	AdvertiserID uuid.UUID
	Name         string
	BudgetTotal  int64
	BudgetDaily  int64
}

type GetCampaignInput struct {
	CampaignID uuid.UUID
}

type PatchCampaignInput struct {
	CampaignID  uuid.UUID
	Name        *string
	BudgetTotal *int64
	BudgetDaily *int64
}

type PutCampaignInput struct {
	CampaignID  uuid.UUID
	Name        string
	BudgetTotal int64
	BudgetDaily int64
}

type DeleteCampaignInput struct {
	CampaignID uuid.UUID
}

type PauseCampaignInput struct {
	CampaignID uuid.UUID
	Reason     string
}

type ResumeCampaignInput struct {
	CampaignID uuid.UUID
}
