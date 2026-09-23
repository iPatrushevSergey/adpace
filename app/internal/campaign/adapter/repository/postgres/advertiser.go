package postgres

import (
	"context"
	"errors"
	"fmt"

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
	transactor *postgres.Transactor
	retryer    port.Retryer
	conv       converter.AdvertiserConverter
}

func NewAdvertiserRepo(transactor *postgres.Transactor, retryer port.Retryer) *AdvertiserRepo {
	return &AdvertiserRepo{
		transactor: transactor,
		retryer:    retryer,
		conv:       &converter.AdvertiserConverterImpl{},
	}
}

func (r *AdvertiserRepo) GetByID(ctx context.Context, advertiserID string) (entity.Advertiser, error) {
	id, err := uuid.Parse(advertiserID)
	if err != nil {
		return entity.Advertiser{}, fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	var m sqlcgen.Advertiser
	err = r.retryer.Do(ctx, func() error {
		q := sqlcgen.New(r.transactor.GetExecutor(ctx))
		var err error
		m, err = q.GetByID(ctx, id)
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

func (r *AdvertiserRepo) GetByIDForUpdate(ctx context.Context, advertiserID string) (entity.Advertiser, error) {
	id, err := uuid.Parse(advertiserID)
	if err != nil {
		return entity.Advertiser{}, fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	var m sqlcgen.Advertiser
	err = r.retryer.Do(ctx, func() error {
		q := sqlcgen.New(r.transactor.GetExecutor(ctx))
		var err error
		m, err = q.GetByIDForUpdate(ctx, id)
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
	params, err := r.conv.ToCreateParams(advertiser)
	if err != nil {
		return entity.Advertiser{}, fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	var m sqlcgen.Advertiser
	err = r.retryer.Do(ctx, func() error {
		q := sqlcgen.New(r.transactor.GetExecutor(ctx))
		var err error
		m, err = q.Create(ctx, params)
		return err
	})
	if err != nil {
		return entity.Advertiser{}, err
	}
	return r.conv.ToEntityAdvertiser(m), nil
}

func (r *AdvertiserRepo) Save(ctx context.Context, advertiser entity.Advertiser) error {
	params, err := r.conv.ToSaveParams(advertiser)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	return r.retryer.Do(ctx, func() error {
		q := sqlcgen.New(r.transactor.GetExecutor(ctx))
		_, err := q.Save(ctx, params)
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	})
}

func (r *AdvertiserRepo) Delete(ctx context.Context, advertiser string) error {
	id, err := uuid.Parse(advertiser)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	return r.retryer.Do(ctx, func() error {
		q := sqlcgen.New(r.transactor.GetExecutor(ctx))
		return q.Delete(ctx, id)
	})
}
