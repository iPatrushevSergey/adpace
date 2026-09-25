package usecase

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
)

type DeleteAdvertiser struct {
	advertiserRepo port.AdvertiserRepo
}

func NewDeleteAdvertiser(advertiserRepo port.AdvertiserRepo) *DeleteAdvertiser {
	return &DeleteAdvertiser{advertiserRepo: advertiserRepo}
}

func (uc *DeleteAdvertiser) Execute(ctx context.Context, in dto.DeleteAdvertiserInput) (out struct{}, err error) {
	if err := uc.advertiserRepo.Delete(ctx, in.AdvertiserID); err != nil {
		return struct{}{}, err
	}
	return struct{}{}, nil
}
