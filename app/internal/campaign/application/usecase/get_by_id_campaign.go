package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type GetByIDCampaign struct {
	campaignRepo port.CampaignRepo
	retryer      port.Retryer
}

func NewGetByIDCampaign(
	campaignRepo port.CampaignRepo,
	retryer port.Retryer,
) *GetByIDCampaign {
	return &GetByIDCampaign{
		campaignRepo: campaignRepo,
		retryer:      retryer,
	}
}

func (uc *GetByIDCampaign) Execute(ctx context.Context, in dto.GetCampaignInput) (out entity.Campaign, err error) {
	err = uc.retryer.Do(ctx, func() error {
		var err error
		out, err = uc.campaignRepo.GetByID(ctx, in.CampaignID)
		return err
	})
	if err != nil {
		return entity.Campaign{}, err
	}
	return out, nil
}
