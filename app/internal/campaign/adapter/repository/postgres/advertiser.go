package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres/converter"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres/sqlcgen"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
	"github.com/iPatrushevSergey/adpace/app/internal/pkg/adapter/repository/postgres"
	"github.com/jackc/pgx/v5"
)

type AdvertiserRepo struct {
	executor *postgres.Executor
	retryer  port.Retryer
	conv     converter.AdvertiserConverter
}

func NewAdvertiserRepo(executor *postgres.Executor, retryer port.Retryer) *AdvertiserRepo {
	return &AdvertiserRepo{
		executor: executor,
		retryer:  retryer,
		conv:     &converter.AdvertiserConverterImpl{},
	}
}

func (r *AdvertiserRepo) GetByID(ctx context.Context, advertiserID uuid.UUID) (entity.Advertiser, error) {
	var m sqlcgen.Advertiser
	err := r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		var err error
		m, err = sqlcgen.New(q).GetAdvertiserByID(ctx, advertiserID)
		return err
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Advertiser{}, domain.ErrNotFound
		}
		return entity.Advertiser{}, err
	}
	return r.conv.ToEntityAdvertiser(m), nil
}

func (r *AdvertiserRepo) GetByIDForUpdate(ctx context.Context, advertiserID uuid.UUID) (entity.Advertiser, error) {
	var m sqlcgen.Advertiser
	err := r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		var err error
		m, err = sqlcgen.New(q).GetAdvertiserByIDForUpdate(ctx, advertiserID)
		return err
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Advertiser{}, domain.ErrNotFound
		}
		return entity.Advertiser{}, err
	}
	return r.conv.ToEntityAdvertiser(m), nil
}

func (r *AdvertiserRepo) Create(ctx context.Context, advertiser entity.Advertiser) (entity.Advertiser, error) {
	params := r.conv.ToCreateAdvertiserParams(advertiser)

	var m sqlcgen.Advertiser
	err := r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		var err error
		m, err = sqlcgen.New(q).CreateAdvertiser(ctx, params)
		return err
	})
	if err != nil {
		return entity.Advertiser{}, err
	}
	return r.conv.ToEntityAdvertiser(m), nil
}

func (r *AdvertiserRepo) Save(ctx context.Context, advertiser entity.Advertiser) error {
	params := r.conv.ToSaveAdvertiserParams(advertiser)

	return r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		_, err := sqlcgen.New(q).SaveAdvertiser(ctx, params)
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	})
}

func (r *AdvertiserRepo) Delete(ctx context.Context, advertiserID uuid.UUID) error {
	return r.executor.Do(ctx, r.retryer, func(q postgres.Querier) error {
		return sqlcgen.New(q).DeleteAdvertiser(ctx, advertiserID)
	})
}
