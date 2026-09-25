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
}

func NewResumeCampaign(
	campaignRepo port.CampaignRepo,
	clock port.Clock,
) *ResumeCampaign {
	return &ResumeCampaign{
		campaignRepo: campaignRepo,
		clock:        clock,
	}
}

func (uc *ResumeCampaign) Execute(ctx context.Context, in dto.ResumeCampaignInput) (out entity.Campaign, err error) {
	return uc.campaignRepo.Resume(ctx, in.CampaignID, uc.clock.Now())
}
