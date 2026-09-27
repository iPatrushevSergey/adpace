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
	retryer      port.Retryer
}

func NewPauseCampaign(
	campaignRepo port.CampaignRepo,
	clock port.Clock,
	retryer port.Retryer,
) *PauseCampaign {
	return &PauseCampaign{
		campaignRepo: campaignRepo,
		clock:        clock,
		retryer:      retryer,
	}
}

func (uc *PauseCampaign) Execute(ctx context.Context, in dto.PauseCampaignInput) (out entity.Campaign, err error) {
	if !entity.IsValidPauseReason(in.Reason) {
		return entity.Campaign{}, fmt.Errorf("%w: invalid pause reason %q", domain.ErrBadInput, in.Reason)
	}

	err = uc.retryer.Do(ctx, func() error {
		var err error
		out, err = uc.campaignRepo.Pause(ctx, in.CampaignID, in.Reason, uc.clock.Now())
		return err
	})
	if err != nil {
		return entity.Campaign{}, err
	}
	return out, nil
}
