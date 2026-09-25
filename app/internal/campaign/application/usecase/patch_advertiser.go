package usecase

import (
	"context"
	"fmt"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type PatchAdvertiser struct {
	advertiserRepo port.AdvertiserRepo
	transactor     port.Transactor
	retryer        port.Retryer
	clock          port.Clock
}

func NewPatchAdvertiser(
	advertiserRepo port.AdvertiserRepo,
	transactor port.Transactor,
	retryer port.Retryer,
	clock port.Clock,
) *PatchAdvertiser {
	return &PatchAdvertiser{
		advertiserRepo: advertiserRepo,
		transactor:     transactor,
		retryer:        retryer,
		clock:          clock,
	}
}

func (uc *PatchAdvertiser) Execute(ctx context.Context, in dto.PatchAdvertiserInput) (out entity.Advertiser, err error) {
	if in.Name == nil && in.Country == nil {
		return entity.Advertiser{}, fmt.Errorf("%w: at least one field must be provided", domain.ErrBadInput)
	}

	err = uc.transactor.RunInTransaction(ctx, uc.retryer, func(ctx context.Context) error {
		current, err := uc.advertiserRepo.GetByIDForUpdate(ctx, in.AdvertiserID)
		if err != nil {
			return err
		}

		if err := current.SetName(in.Name); err != nil {
			return err
		}
		if err := current.SetCountry(in.Country); err != nil {
			return err
		}
		current.UpdatedAt = uc.clock.Now()

		if err := uc.advertiserRepo.Save(ctx, current); err != nil {
			return err
		}
		out = current
		return nil
	})
	if err != nil {
		return entity.Advertiser{}, err
	}
	return out, nil
}
