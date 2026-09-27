package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type GetByIDAdvertiser struct {
	advertiserRepo port.AdvertiserRepo
	retryer        port.Retryer
}

func NewGetByIDAdvertiser(
	advertiserRepo port.AdvertiserRepo,
	retryer port.Retryer,
) *GetByIDAdvertiser {
	return &GetByIDAdvertiser{
		advertiserRepo: advertiserRepo,
		retryer:        retryer,
	}
}

func (uc *GetByIDAdvertiser) Execute(ctx context.Context, in dto.GetAdvertiserInput) (out entity.Advertiser, err error) {
	err = uc.retryer.Do(ctx, func() error {
		var err error
		out, err = uc.advertiserRepo.GetByID(ctx, in.AdvertiserID)
		return err
	})
	if err != nil {
		return entity.Advertiser{}, err
	}
	return out, nil
}
