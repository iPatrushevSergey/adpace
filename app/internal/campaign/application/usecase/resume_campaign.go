package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type ResumeCampaign struct {
	campaignRepo port.CampaignRepo
	clock        port.Clock
	retryer      port.Retryer
}

func NewResumeCampaign(
	campaignRepo port.CampaignRepo,
	clock port.Clock,
	retryer port.Retryer,
) *ResumeCampaign {
	return &ResumeCampaign{
		campaignRepo: campaignRepo,
		clock:        clock,
		retryer:      retryer,
	}
}

func (uc *ResumeCampaign) Execute(ctx context.Context, in dto.ResumeCampaignInput) (out entity.Campaign, err error) {
	err = uc.retryer.Do(ctx, func() error {
		var err error
		out, err = uc.campaignRepo.Resume(ctx, in.CampaignID, uc.clock.Now())
		return err
	})
	if err != nil {
		return entity.Campaign{}, err
	}
	return out, nil
}
