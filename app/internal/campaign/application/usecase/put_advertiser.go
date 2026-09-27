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
	retryer        port.Retryer
}

func NewPutAdvertiser(
	advertiserRepo port.AdvertiserRepo,
	clock port.Clock,
	retryer port.Retryer,
) *PutAdvertiser {
	return &PutAdvertiser{
		advertiserRepo: advertiserRepo,
		clock:          clock,
		retryer:        retryer,
	}
}

func (uc *PutAdvertiser) Execute(ctx context.Context, in dto.PutAdvertiserInput) (out struct{}, err error) {
	advertiser, err := entity.NewAdvertiser(in.AdvertiserID, in.Name, in.Country, uc.clock.Now())
	if err != nil {
		return struct{}{}, err
	}

	err = uc.retryer.Do(ctx, func() error {
		return uc.advertiserRepo.Save(ctx, advertiser)
	})
	if err != nil {
		return struct{}{}, err
	}
	return struct{}{}, nil
}
