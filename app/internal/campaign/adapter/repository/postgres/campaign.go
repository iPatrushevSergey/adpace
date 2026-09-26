package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres/converter"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres/sqlcgen"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
	"github.com/iPatrushevSergey/adpace/app/internal/pkg/adapter/repository/postgres"
	"github.com/jackc/pgx/v5"
)

type CampaignRepo struct {
	executor *postgres.Executor
	retryer  port.Retryer
	conv     converter.CampaignConverter
}

func NewCampaignRepo(executor *postgres.Executor, retryer port.Retryer) *CampaignRepo {
	return &CampaignRepo{
		executor: executor,
		retryer:  retryer,
		conv:     &converter.CampaignConverterImpl{},
	}
}

func (r *CampaignRepo) GetByID(ctx context.Context, campaignID uuid.UUID) (entity.Campaign, error) {
	var m sqlcgen.Campaign
	err := r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		var err error
		m, err = sqlcgen.New(q).GetCampaignByID(ctx, campaignID)
		return err
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Campaign{}, domain.ErrNotFound
		}
		return entity.Campaign{}, err
	}
	return r.conv.ToEntityCampaign(m), nil
}

func (r *CampaignRepo) GetByIDForUpdate(ctx context.Context, campaignID uuid.UUID) (entity.Campaign, error) {
	var m sqlcgen.Campaign
	err := r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		var err error
		m, err = sqlcgen.New(q).GetCampaignByIDForUpdate(ctx, campaignID)
		return err
	})
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

	var m sqlcgen.Campaign
	err := r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		var err error
		m, err = sqlcgen.New(q).CreateCampaign(ctx, params)
		return err
	})
	if err != nil {
		return entity.Campaign{}, err
	}
	return r.conv.ToEntityCampaign(m), nil
}

func (r *CampaignRepo) Save(ctx context.Context, campaign entity.Campaign) error {
	params := r.conv.ToSaveCampaignParams(campaign)

	return r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		_, err := sqlcgen.New(q).SaveCampaign(ctx, params)
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	})
}

func (r *CampaignRepo) Pause(ctx context.Context, campaignID uuid.UUID, reason string, updatedAt time.Time) (entity.Campaign, error) {
	var m sqlcgen.Campaign
	err := r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		var err error
		m, err = sqlcgen.New(q).PauseCampaign(ctx, sqlcgen.PauseCampaignParams{
			CampaignID:  campaignID,
			PauseReason: &reason,
			UpdatedAt:   updatedAt,
		})
		return err
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
	var m sqlcgen.Campaign
	err := r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		var err error
		m, err = sqlcgen.New(q).ResumeCampaign(ctx, sqlcgen.ResumeCampaignParams{
			CampaignID: campaignID,
			UpdatedAt:  updatedAt,
		})
		return err
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
	return r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		return sqlcgen.New(q).DeleteCampaign(ctx, campaignID)
	})
}
