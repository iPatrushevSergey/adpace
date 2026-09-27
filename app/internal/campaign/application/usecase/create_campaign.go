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
	retryer      port.Retryer
}

func NewCreateCampaign(
	campaignRepo port.CampaignRepo,
	idGenerator port.IDGenerator,
	clock port.Clock,
	retryer port.Retryer,
) *CreateCampaign {
	return &CreateCampaign{
		campaignRepo: campaignRepo,
		idGenerator:  idGenerator,
		clock:        clock,
		retryer:      retryer,
	}
}

func (uc *CreateCampaign) Execute(ctx context.Context, in dto.CreateCampaignInput) (out entity.Campaign, err error) {
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

	err = uc.retryer.Do(ctx, func() error {
		var err error
		out, err = uc.campaignRepo.Create(ctx, campaign)
		return err
	})
	if err != nil {
		return entity.Campaign{}, err
	}
	return out, nil
}
