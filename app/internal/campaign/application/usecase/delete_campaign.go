package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
)

type DeleteCampaign struct {
	campaignRepo port.CampaignRepo
	retryer      port.Retryer
}

func NewDeleteCampaign(
	campaignRepo port.CampaignRepo,
	retryer port.Retryer,
) *DeleteCampaign {
	return &DeleteCampaign{campaignRepo: campaignRepo, retryer: retryer}
}

func (uc *DeleteCampaign) Execute(ctx context.Context, in dto.DeleteCampaignInput) (struct{}, error) {
	err := uc.retryer.Do(ctx, func() error {
		return uc.campaignRepo.Delete(ctx, in.CampaignID)
	})
	if err != nil {
		return struct{}{}, err
	}
	return struct{}{}, nil
}
