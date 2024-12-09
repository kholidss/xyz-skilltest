package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kholidss/xyz-skilltest/internal/entity"
	"github.com/kholidss/xyz-skilltest/pkg/database/mysql"

	"github.com/kholidss/xyz-skilltest/pkg/helper"
	"github.com/kholidss/xyz-skilltest/pkg/tracer"
	"github.com/pkg/errors"
)

type bucketRepository struct {
	db mysql.Adapter
}

func NewBucketRepository(db mysql.Adapter) BucketRepository {
	return &bucketRepository{
		db: db,
	}
}

func (b bucketRepository) Store(ctx context.Context, payload any, opts ...Option) error {
	var (
		err error
		tx  *sql.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "BucketRepo.Store", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = b.db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		if err != nil {
			tracer.AddSpanError(span, err)
			return err
		}

		defer func() {
			err = tx.Commit()
			if err != nil {
				tracer.AddSpanError(span, err)
				err = errors.Wrap(err, "failed to commit")
			}
		}()
	}

	query, val, err := helper.StructQueryInsertMysql(payload, TableNameBuckets, "db", false)

	_, err = tx.ExecContext(
		ctx,
		query,
		val...,
	)
	if err != nil {
		tracer.AddSpanError(span, err)
		return err
	}

	return err
}

func (b bucketRepository) Update(ctx context.Context, payload any, where any, opts ...Option) error {
	var (
		err error
		tx  *sql.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "BucketRepo.Update", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = b.db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		if err != nil {
			tracer.AddSpanError(span, err)
			return err
		}

		defer func() {
			err = tx.Commit()
			if err != nil {
				err = errors.Wrap(err, "failed to commit")
			}
		}()
	}

	q, vals, err := helper.StructToQueryUpdateMysql(payload, where, TableNameBuckets, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return err
	}

	_, err = tx.ExecContext(ctx, q, vals...)
	if err != nil {
		tracer.AddSpanError(span, err)
		return err
	}

	return err
}

func (b *bucketRepository) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.Bucket, error) {
	var (
		dest entity.Bucket
	)

	ctx, span := tracer.NewSpan(ctx, "BucketRepo.FindOne", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	err = b.db.QueryRow(ctx, &dest, fmt.Sprintf(DefaultQueryFindOne, helper.SelectCustom(selectColumn), TableNameBuckets, wq), vals...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return &dest, nil
}

func (b *bucketRepository) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Bucket, error) {
	var (
		dest []entity.Bucket
	)

	ctx, span := tracer.NewSpan(ctx, "BucketRepo.Finds", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	err = b.db.Query(ctx, &dest, fmt.Sprintf(DefaultQueryFinds, helper.SelectCustom(selectColumns), TableNameBuckets, wq), vals...)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return dest, nil
}

func (b bucketRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return b.db.BeginTx(ctx, opts)
}
