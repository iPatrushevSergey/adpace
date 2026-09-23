package port

import (
	"context"
	"time"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

type CampaignRepo interface {
	GetByID(ctx context.Context, campaignID string) (entity.Campaign, error)
	GetByIDForUpdate(ctx context.Context, campaignID string) (entity.Campaign, error)
	Create(ctx context.Context, campaign entity.Campaign) (entity.Campaign, error)
	Save(ctx context.Context, campaign entity.Campaign) error
	Pause(ctx context.Context, campaignID, reason string, updatedAt time.Time) (entity.Campaign, error)
	Resume(ctx context.Context, campaignID string, updatedAt time.Time) (entity.Campaign, error)
	Delete(ctx context.Context, campaignID string) error
}
