package usecase

import (
	"context"
	"fmt"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type PauseCampaign struct {
	campaignRepo port.CampaignRepo
	clock        port.Clock
}

func NewPauseCampaign(
	campaignRepo port.CampaignRepo,
	clock port.Clock,
) *PauseCampaign {
	return &PauseCampaign{
		campaignRepo: campaignRepo,
		clock:        clock,
	}
}

func (uc *PauseCampaign) Execute(ctx context.Context, in dto.PauseCampaignInput) (out entity.Campaign, err error) {
	if !entity.IsValidPauseReason(in.Reason) {
		return entity.Campaign{}, fmt.Errorf("%w: invalid pause reason %q", domain.ErrBadInput, in.Reason)
	}

	return uc.campaignRepo.Pause(ctx, in.CampaignID, in.Reason, uc.clock.Now())
}
