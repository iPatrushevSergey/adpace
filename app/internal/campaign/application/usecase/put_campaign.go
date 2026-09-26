package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type PutCampaign struct {
	campaignRepo port.CampaignRepo
	clock        port.Clock
}

func NewPutCampaign(campaignRepo port.CampaignRepo, clock port.Clock) *PutCampaign {
	return &PutCampaign{campaignRepo: campaignRepo, clock: clock}
}

func (uc *PutCampaign) Execute(ctx context.Context, in dto.PutCampaignInput) (struct{}, error) {
	// AdvertiserID intentionally not part of Put — ownership is immutable after Create.
	campaign, err := entity.NewCampaign(
		in.CampaignID,
		uuid.Nil,
		in.Name,
		in.BudgetTotal,
		in.BudgetDaily,
		uc.clock.Now(),
	)
	if err != nil {
		return struct{}{}, err
	}

	if err := uc.campaignRepo.Save(ctx, campaign); err != nil {
		return struct{}{}, err
	}
	return struct{}{}, nil
}
