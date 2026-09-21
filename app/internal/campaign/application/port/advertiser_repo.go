package port

import (
	"context"
	"time"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type AdvertiserRepo interface {
	GetByID(ctx context.Context, advertiserID string) (entity.Advertiser, error)
	Create(ctx context.Context, advertiser entity.Advertiser) (entity.Advertiser, error)
	Patch(ctx context.Context, patch dto.PatchAdvertiserInput, updatedAt time.Time) (entity.Advertiser, error)
	Put(ctx context.Context, advertiser entity.Advertiser) error
	Delete(ctx context.Context, advertiserID string) error
}
