package usecase

import (
	"context"
	"fmt"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
	"github.com/iPatrushevSergey/adpace/app/internal/pkg/apputil"
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
	if !apputil.IsUUID(in.CampaignID) {
		return entity.Campaign{}, fmt.Errorf("%w: campaign id is not a valid UUID", domain.ErrBadInput)
	}

	return uc.campaignRepo.Resume(ctx, in.CampaignID, uc.clock.Now())
}
