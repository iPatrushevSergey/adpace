package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type AdvertiserRepo interface {
	GetByID(ctx context.Context, advertiserID uuid.UUID) (entity.Advertiser, error)
	GetByIDForUpdate(ctx context.Context, advertiserID uuid.UUID) (entity.Advertiser, error)
	Create(ctx context.Context, advertiser entity.Advertiser) (entity.Advertiser, error)
	Save(ctx context.Context, advertiser entity.Advertiser) error
	Delete(ctx context.Context, advertiserID uuid.UUID) error
}
