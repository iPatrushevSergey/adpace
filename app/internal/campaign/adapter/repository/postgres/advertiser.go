package postgres

import (
	"context"
	"errors"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres/converter"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres/sqlcgen"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdvertiserRepo struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	conv   converter.AdvertiserConverter
}

func NewAdvertiserRepo(pool *pgxpool.Pool, getter *trmpgx.CtxGetter) *AdvertiserRepo {
	return &AdvertiserRepo{
		pool:   pool,
		getter: getter,
		conv:   &converter.AdvertiserConverterImpl{},
	}
}

func (r *AdvertiserRepo) conn(ctx context.Context) trmpgx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.pool)
}

func (r *AdvertiserRepo) GetByID(ctx context.Context, advertiserID uuid.UUID) (entity.Advertiser, error) {
	m, err := sqlcgen.New(r.conn(ctx)).GetAdvertiserByID(ctx, advertiserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Advertiser{}, domain.ErrNotFound
		}
		return entity.Advertiser{}, err
	}
	return r.conv.ToEntityAdvertiser(m), nil
}

func (r *AdvertiserRepo) GetByIDForUpdate(ctx context.Context, advertiserID uuid.UUID) (entity.Advertiser, error) {
	m, err := sqlcgen.New(r.conn(ctx)).GetAdvertiserByIDForUpdate(ctx, advertiserID)
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

	m, err := sqlcgen.New(r.conn(ctx)).CreateAdvertiser(ctx, params)
	if err != nil {
		return entity.Advertiser{}, err
	}
	return r.conv.ToEntityAdvertiser(m), nil
}

func (r *AdvertiserRepo) Save(ctx context.Context, advertiser entity.Advertiser) error {
	params := r.conv.ToSaveAdvertiserParams(advertiser)

	_, err := sqlcgen.New(r.conn(ctx)).SaveAdvertiser(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func (r *AdvertiserRepo) Delete(ctx context.Context, advertiserID uuid.UUID) error {
	return sqlcgen.New(r.conn(ctx)).DeleteAdvertiser(ctx, advertiserID)
}
