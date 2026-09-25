package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type PutAdvertiser struct {
	advertiserRepo port.AdvertiserRepo
	clock          port.Clock
}

func NewPutAdvertiser(
	advertiserRepo port.AdvertiserRepo,
	clock port.Clock,
) *PutAdvertiser {
	return &PutAdvertiser{
		advertiserRepo: advertiserRepo,
		clock:          clock,
	}
}

func (uc *PutAdvertiser) Execute(ctx context.Context, in dto.PutAdvertiserInput) (out struct{}, err error) {
	advertiser, err := entity.NewAdvertiser(in.AdvertiserID, in.Name, in.Country, uc.clock.Now())
	if err != nil {
		return struct{}{}, err
	}

	if err := uc.advertiserRepo.Save(ctx, advertiser); err != nil {
		return struct{}{}, err
	}
	return struct{}{}, nil
}
