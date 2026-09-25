package postgres

import (
	"context"
	"errors"
	"fmt"
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

func (r *CampaignRepo) GetByID(ctx context.Context, campaignID string) (entity.Campaign, error) {
	id, err := uuid.Parse(campaignID)
	if err != nil {
		return entity.Campaign{}, fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	var m sqlcgen.Campaign
	err = r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		var err error
		m, err = sqlcgen.New(q).GetCampaignByID(ctx, id)
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

func (r *CampaignRepo) GetByIDForUpdate(ctx context.Context, campaignID string) (entity.Campaign, error) {
	id, err := uuid.Parse(campaignID)
	if err != nil {
		return entity.Campaign{}, fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	var m sqlcgen.Campaign
	err = r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		var err error
		m, err = sqlcgen.New(q).GetCampaignByIDForUpdate(ctx, id)
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
	params, err := r.conv.ToCreateCampaignParams(campaign)
	if err != nil {
		return entity.Campaign{}, fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	var m sqlcgen.Campaign
	err = r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
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
	params, err := r.conv.ToSaveCampaignParams(campaign)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	return r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		_, err := sqlcgen.New(q).SaveCampaign(ctx, params)
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	})
}

func (r *CampaignRepo) Pause(ctx context.Context, campaignID, reason string, updatedAt time.Time) (entity.Campaign, error) {
	id, err := uuid.Parse(campaignID)
	if err != nil {
		return entity.Campaign{}, fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	var m sqlcgen.Campaign
	err = r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		var err error
		m, err = sqlcgen.New(q).PauseCampaign(ctx, sqlcgen.PauseCampaignParams{
			CampaignID:  id,
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

func (r *CampaignRepo) Resume(ctx context.Context, campaignID string, updatedAt time.Time) (entity.Campaign, error) {
	id, err := uuid.Parse(campaignID)
	if err != nil {
		return entity.Campaign{}, fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	var m sqlcgen.Campaign
	err = r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		var err error
		m, err = sqlcgen.New(q).ResumeCampaign(ctx, sqlcgen.ResumeCampaignParams{
			CampaignID: id,
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

func (r *CampaignRepo) Delete(ctx context.Context, campaignID string) error {
	id, err := uuid.Parse(campaignID)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	return r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		return sqlcgen.New(q).DeleteCampaign(ctx, id)
	})
}
