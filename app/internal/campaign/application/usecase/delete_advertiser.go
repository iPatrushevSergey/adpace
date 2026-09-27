package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
)

type DeleteAdvertiser struct {
	advertiserRepo port.AdvertiserRepo
	retryer        port.Retryer
}

func NewDeleteAdvertiser(
	advertiserRepo port.AdvertiserRepo,
	retryer port.Retryer,
) *DeleteAdvertiser {
	return &DeleteAdvertiser{advertiserRepo: advertiserRepo, retryer: retryer}
}

func (uc *DeleteAdvertiser) Execute(ctx context.Context, in dto.DeleteAdvertiserInput) (out struct{}, err error) {
	err = uc.retryer.Do(ctx, func() error {
		return uc.advertiserRepo.Delete(ctx, in.AdvertiserID)
	})
	if err != nil {
		return struct{}{}, err
	}
	return struct{}{}, nil
}
