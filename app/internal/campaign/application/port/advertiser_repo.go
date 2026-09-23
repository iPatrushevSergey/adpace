package port

import (
	"context"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type AdvertiserRepo interface {
	GetByID(ctx context.Context, advertiserID string) (entity.Advertiser, error)
	GetByIDForUpdate(ctx context.Context, advertiserID string) (entity.Advertiser, error)
	Create(ctx context.Context, advertiser entity.Advertiser) (entity.Advertiser, error)
	Save(ctx context.Context, advertiser entity.Advertiser) error
	Delete(ctx context.Context, advertiserID string) error
}
