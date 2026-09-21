package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres/converter"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres/sqlcgen"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
	"github.com/iPatrushevSergey/adpace/app/internal/pkg/adapter/repository/postgres"
	"github.com/iPatrushevSergey/adpace/app/internal/pkg/apputil"
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

func (r *AdvertiserRepo) Put(ctx context.Context, advertiser entity.Advertiser) error {
	params, err := r.conv.ToPutParams(advertiser)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	return r.retryer.Do(ctx, func() error {
		q := sqlcgen.New(r.transactor.GetExecutor(ctx))
		_, err := q.Put(ctx, params)
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	})
}

func (r *AdvertiserRepo) Patch(ctx context.Context, patch dto.PatchAdvertiserInput, updatedAt time.Time) (entity.Advertiser, error) {
	id, err := uuid.Parse(patch.AdvertiserID)
	if err != nil {
		return entity.Advertiser{}, fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	b := apputil.NewPatchQueryBuilder(
		"advertiser", "advertiser_id", id,
		"advertiser_id, name, country, created_at, updated_at",
		3,
	)
	apputil.Set(b, "name", patch.Name)
	apputil.Set(b, "country", patch.Country)
	apputil.Set(b, "updated_at", &updatedAt)

	query, args, err := b.Build()
	if err != nil {
		return entity.Advertiser{}, fmt.Errorf("%w: %v", domain.ErrBadInput, err)
	}

	var m sqlcgen.Advertiser
	err = r.retryer.Do(ctx, func() error {
		q := r.transactor.GetExecutor(ctx)
		row := q.QueryRow(ctx, query, args...)
		return row.Scan(&m.AdvertiserID, &m.Name, &m.Country, &m.CreatedAt, &m.UpdatedAt)
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Advertiser{}, domain.ErrNotFound
		}
		return entity.Advertiser{}, err
	}
	return r.conv.ToEntityAdvertiser(m), nil
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
