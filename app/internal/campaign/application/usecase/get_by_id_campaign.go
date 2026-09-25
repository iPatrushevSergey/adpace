package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type GetByIDCampaign struct {
	campaignRepo port.CampaignRepo
}

func NewGetByIDCampaign(campaignRepo port.CampaignRepo) *GetByIDCampaign {
	return &GetByIDCampaign{campaignRepo: campaignRepo}
}

func (uc *GetByIDCampaign) Execute(ctx context.Context, in dto.GetCampaignInput) (entity.Campaign, error) {
	return uc.campaignRepo.GetByID(ctx, in.CampaignID)
}
