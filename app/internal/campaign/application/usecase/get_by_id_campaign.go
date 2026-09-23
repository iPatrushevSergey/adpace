package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
	"github.com/iPatrushevSergey/adpace/app/internal/pkg/apputil"
)

type GetByIDCampaign struct {
	campaignRepo port.CampaignRepo
}

func NewGetByIDCampaign(campaignRepo port.CampaignRepo) *GetByIDCampaign {
	return &GetByIDCampaign{campaignRepo: campaignRepo}
}

func (uc *GetByIDCampaign) Execute(ctx context.Context, in dto.GetCampaignInput) (entity.Campaign, error) {
	if !apputil.IsUUID(in.CampaignID) {
		return entity.Campaign{}, domain.ErrBadInput
	}
	return uc.campaignRepo.GetByID(ctx, in.CampaignID)
}
