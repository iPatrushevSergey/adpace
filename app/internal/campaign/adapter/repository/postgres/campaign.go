package postgres

import (
	"context"
	"errors"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres/converter"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres/sqlcgen"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CampaignRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	conv   converter.CampaignConverter
}

func NewCampaignRepo(pool *pgxpool.Pool, getter *trmpgx.CtxGetter) *CampaignRepo {
	return &CampaignRepo{
		pool:   pool,
		getter: getter,
		conv:   &converter.CampaignConverterImpl{},
	}
}

func (r *CampaignRepo) conn(ctx context.Context) trmpgx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.pool)
}

func (r *CampaignRepo) GetByID(ctx context.Context, campaignID uuid.UUID) (entity.Campaign, error) {
	m, err := sqlcgen.New(r.conn(ctx)).GetCampaignByID(ctx, campaignID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Campaign{}, domain.ErrNotFound
		}
		return entity.Campaign{}, err
	}
	return r.conv.ToEntityCampaign(m), nil
}

func (r *CampaignRepo) GetByIDForUpdate(ctx context.Context, campaignID uuid.UUID) (entity.Campaign, error) {
	m, err := sqlcgen.New(r.conn(ctx)).GetCampaignByIDForUpdate(ctx, campaignID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Campaign{}, domain.ErrNotFound
		}
		return entity.Campaign{}, err
	}
	return r.conv.ToEntityCampaign(m), nil
}

func (r *CampaignRepo) Create(ctx context.Context, campaign entity.Campaign) (entity.Campaign, error) {
	params := r.conv.ToCreateCampaignParams(campaign)

	m, err := sqlcgen.New(r.conn(ctx)).CreateCampaign(ctx, params)
	if err != nil {
		return entity.Campaign{}, err
	}
	return r.conv.ToEntityCampaign(m), nil
}

func (r *CampaignRepo) Save(ctx context.Context, campaign entity.Campaign) error {
	params := r.conv.ToSaveCampaignParams(campaign)

	_, err := sqlcgen.New(r.conn(ctx)).SaveCampaign(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func (r *CampaignRepo) Pause(ctx context.Context, campaignID uuid.UUID, reason string, updatedAt time.Time) (entity.Campaign, error) {
	m, err := sqlcgen.New(r.conn(ctx)).PauseCampaign(ctx, sqlcgen.PauseCampaignParams{
		CampaignID:  campaignID,
		PauseReason: &reason,
		UpdatedAt:   updatedAt,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Campaign{}, domain.ErrConflict
		}
		return entity.Campaign{}, err
	}
	return r.conv.ToEntityCampaign(m), nil
}

func (r *CampaignRepo) Resume(ctx context.Context, campaignID uuid.UUID, updatedAt time.Time) (entity.Campaign, error) {
	m, err := sqlcgen.New(r.conn(ctx)).ResumeCampaign(ctx, sqlcgen.ResumeCampaignParams{
		CampaignID: campaignID,
		UpdatedAt:  updatedAt,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Campaign{}, domain.ErrConflict
		}
		return entity.Campaign{}, err
	}
	return r.conv.ToEntityCampaign(m), nil
}

func (r *CampaignRepo) Delete(ctx context.Context, campaignID uuid.UUID) error {
	return sqlcgen.New(r.conn(ctx)).DeleteCampaign(ctx, campaignID)
}
