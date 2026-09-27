package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type CreateAdvertiser struct {
	advertiserRepo port.AdvertiserRepo
	idGenerator    port.IDGenerator
	clock          port.Clock
	retryer        port.Retryer
}

func NewCreateAdvertiser(
	advertiserRepo port.AdvertiserRepo,
	idGenerator port.IDGenerator,
	clock port.Clock,
	retryer port.Retryer,
) *CreateAdvertiser {
	return &CreateAdvertiser{
		advertiserRepo: advertiserRepo,
		idGenerator:    idGenerator,
		clock:          clock,
		retryer:        retryer,
	}
}

func (uc *CreateAdvertiser) Execute(ctx context.Context, in dto.CreateAdvertiserInput) (out entity.Advertiser, err error) {
	id, err := uc.idGenerator.NewID()
	if err != nil {
		return entity.Advertiser{}, err
	}
	now := uc.clock.Now()

	advertiser, err := entity.NewAdvertiser(id, in.Name, in.Country, now, entity.WithAdvertiserCreatedAt(now))
	if err != nil {
		return entity.Advertiser{}, err
	}

	err = uc.retryer.Do(ctx, func() error {
		var err error
		out, err = uc.advertiserRepo.Create(ctx, advertiser)
		return err
	})
	if err != nil {
		return entity.Advertiser{}, err
	}
	return out, nil
}
