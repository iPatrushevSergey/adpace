package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type CreateCampaign struct {
	campaignRepo port.CampaignRepo
	idGenerator  port.IDGenerator
	clock        port.Clock
}

func NewCreateCampaign(
	campaignRepo port.CampaignRepo,
	idGenerator port.IDGenerator,
	clock port.Clock,
) *CreateCampaign {
	return &CreateCampaign{
		campaignRepo: campaignRepo,
		idGenerator:  idGenerator,
		clock:        clock,
	}
}

func (uc *CreateCampaign) Execute(ctx context.Context, in dto.CreateCampaignInput) (entity.Campaign, error) {
	id, err := uc.idGenerator.NewID()
	if err != nil {
		return entity.Campaign{}, err
	}
	now := uc.clock.Now()

	campaign, err := entity.NewCampaign(
		id,
		in.AdvertiserID,
		in.Name,
		in.BudgetTotal,
		in.BudgetDaily,
		now,
		entity.WithCampaignCreatedAt(now),
	)
	if err != nil {
		return entity.Campaign{}, err
	}

	created, err := uc.campaignRepo.Create(ctx, campaign)
	if err != nil {
		return entity.Campaign{}, err
	}
	return created, nil
}
