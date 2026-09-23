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

type PatchCampaign struct {
	campaignRepo port.CampaignRepo
	transactor   port.Transactor
	retryer      port.Retryer
	clock        port.Clock
}

func NewPatchCampaign(
	campaignRepo port.CampaignRepo,
	transactor port.Transactor,
	retryer port.Retryer,
	clock port.Clock,
) *PatchCampaign {
	return &PatchCampaign{
		campaignRepo: campaignRepo,
		transactor:   transactor,
		retryer:      retryer,
		clock:        clock,
	}
}

func (uc *PatchCampaign) Execute(ctx context.Context, in dto.PatchCampaignInput) (out entity.Campaign, err error) {
	if !apputil.IsUUID(in.CampaignID) {
		return entity.Campaign{}, fmt.Errorf("%w: campaign id is not a valid UUID", domain.ErrBadInput)
	}
	if in.Name == nil && in.BudgetTotal == nil && in.BudgetDaily == nil {
		return entity.Campaign{}, fmt.Errorf("%w: at least one field must be provided", domain.ErrBadInput)
	}

	err = uc.transactor.RunInTransaction(ctx, uc.retryer, func(ctx context.Context) error {
		current, err := uc.campaignRepo.GetByIDForUpdate(ctx, in.CampaignID)
		if err != nil {
			return err
		}

		if err := current.SetName(in.Name); err != nil {
			return err
		}
		if err := current.SetBudgets(in.BudgetTotal, in.BudgetDaily); err != nil {
			return err
		}
		current.UpdatedAt = uc.clock.Now()

		if err := uc.campaignRepo.Save(ctx, current); err != nil {
			return err
		}
		out = current
		return nil
	})
	if err != nil {
		return entity.Campaign{}, err
	}
	return out, nil
}
