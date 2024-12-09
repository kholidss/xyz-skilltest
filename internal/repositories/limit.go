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

type limitRepository struct {
	db mysql.Adapter
}

func NewLimitRepository(db mysql.Adapter) LimitRepository {
	return &limitRepository{
		db: db,
	}
}

func (l limitRepository) Store(ctx context.Context, payload any, opts ...Option) error {
	var (
		err error
		tx  *sql.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "LimitRepo.Store", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = l.db.BeginTx(ctx, &sql.TxOptions{
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

	query, val, err := helper.StructQueryInsertMysql(payload, TableNameLimits, "db", false)

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

func (l limitRepository) Update(ctx context.Context, payload any, where any, opts ...Option) error {
	var (
		err error
		tx  *sql.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "LimitRepo.Update", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = l.db.BeginTx(ctx, &sql.TxOptions{
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

	q, vals, err := helper.StructToQueryUpdateMysql(payload, where, TableNameLimits, "db")
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

func (l *limitRepository) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.Limit, error) {
	var (
		dest entity.Limit
	)

	ctx, span := tracer.NewSpan(ctx, "LimitRepo.FindOne", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	err = l.db.QueryRow(ctx, &dest, fmt.Sprintf(DefaultQueryFindOne, helper.SelectCustom(selectColumn), TableNameLimits, wq), vals...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return &dest, nil
}

func (l *limitRepository) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Limit, error) {
	var (
		dest []entity.Limit
	)

	ctx, span := tracer.NewSpan(ctx, "LimitRepo.Finds", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	err = l.db.Query(ctx, &dest, fmt.Sprintf(DefaultQueryFinds, helper.SelectCustom(selectColumns), TableNameLimits, wq), vals...)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return dest, nil
}

func (l limitRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return l.db.BeginTx(ctx, opts)
}
