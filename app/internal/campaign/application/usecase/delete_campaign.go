package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
)

type DeleteCampaign struct {
	campaignRepo port.CampaignRepo
}

func NewDeleteCampaign(campaignRepo port.CampaignRepo) *DeleteCampaign {
	return &DeleteCampaign{campaignRepo: campaignRepo}
}

func (uc *DeleteCampaign) Execute(ctx context.Context, in dto.DeleteCampaignInput) (struct{}, error) {
	if err := uc.campaignRepo.Delete(ctx, in.CampaignID); err != nil {
		return struct{}{}, err
	}
	return struct{}{}, nil
}
