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
}

func NewCreateAdvertiser(
	advertiserRepo port.AdvertiserRepo,
	idGenerator port.IDGenerator,
	clock port.Clock,
) *CreateAdvertiser {
	return &CreateAdvertiser{
		advertiserRepo: advertiserRepo,
		idGenerator:    idGenerator,
		clock:          clock,
	}
}

func (uc *CreateAdvertiser) Execute(ctx context.Context, in dto.CreateAdvertiserInput) (out entity.Advertiser, err error) {
	id, err := uc.idGenerator.NewID()
	if err != nil {
		return entity.Advertiser{}, err
	}
	now := uc.clock.Now()

	advertiser, err := entity.NewAdvertiser(id, in.Name, in.Country, now, entity.WithCreatedAt(now))
	if err != nil {
		return entity.Advertiser{}, err
	}

	created, err := uc.advertiserRepo.Create(ctx, advertiser)
	if err != nil {
		return entity.Advertiser{}, err
	}
	return created, nil
}
