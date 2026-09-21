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

type PatchAdvertiser struct {
	advertiserRepo port.AdvertiserRepo
	clock          port.Clock
}

func NewPatchAdvertiser(
	advertiserRepo port.AdvertiserRepo,
	clock port.Clock,
) *PatchAdvertiser {
	return &PatchAdvertiser{
		advertiserRepo: advertiserRepo,
		clock:          clock,
	}
}

func (uc *PatchAdvertiser) Execute(ctx context.Context, in dto.PatchAdvertiserInput) (out entity.Advertiser, err error) {
	if !apputil.IsUUID(in.AdvertiserID) {
		return entity.Advertiser{}, fmt.Errorf("%w: advertiser id is not a valid UUID", domain.ErrBadInput)
	}
	if in.Name == nil && in.Country == nil {
		return entity.Advertiser{}, fmt.Errorf("%w: at least one field must be provided", domain.ErrBadInput)
	}
	if in.Name != nil && *in.Name == "" {
		return entity.Advertiser{}, fmt.Errorf("%w: name cannot be empty", domain.ErrBadInput)
	}
	if in.Country != nil && !entity.IsValidCountry(*in.Country) {
		return entity.Advertiser{}, fmt.Errorf("%w: there is no such country code", domain.ErrBadInput)
	}

	updated, err := uc.advertiserRepo.Patch(ctx, in, uc.clock.Now())
	if err != nil {
		return entity.Advertiser{}, err
	}
	return updated, nil
}
