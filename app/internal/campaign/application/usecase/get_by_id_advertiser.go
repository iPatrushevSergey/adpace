package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type GetByIDAdvertiser struct {
	advertiserRepo port.AdvertiserRepo
}

func NewGetByIDAdvertiser(
	advertiserRepo port.AdvertiserRepo,
) *GetByIDAdvertiser {
	return &GetByIDAdvertiser{
		advertiserRepo: advertiserRepo,
	}
}

func (uc *GetByIDAdvertiser) Execute(ctx context.Context, in dto.GetAdvertiserInput) (out entity.Advertiser, err error) {
	advertiser, err := uc.advertiserRepo.GetByID(ctx, in.AdvertiserID)
	if err != nil {
		return entity.Advertiser{}, err
	}
	return advertiser, nil
}
