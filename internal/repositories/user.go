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

type userRepository struct {
	db mysql.Adapter
}

func NewUserRepository(db mysql.Adapter) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u userRepository) Store(ctx context.Context, payload any, opts ...Option) error {
	var (
		err error
		tx  *sql.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "UserRepo.Store", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = u.db.BeginTx(ctx, &sql.TxOptions{
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

	query, val, err := helper.StructQueryInsertMysql(payload, TableNameUsers, "db", false)

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

func (u userRepository) Update(ctx context.Context, payload any, where any, opts ...Option) error {
	var (
		err error
		tx  *sql.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "UserRepo.Update", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = u.db.BeginTx(ctx, &sql.TxOptions{
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

	q, vals, err := helper.StructToQueryUpdateMysql(payload, where, TableNameUsers, "db")
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

func (u *userRepository) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.User, error) {
	var (
		dest entity.User
	)

	ctx, span := tracer.NewSpan(ctx, "UserRepo.FindOne", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	err = u.db.QueryRow(ctx, &dest, fmt.Sprintf(DefaultQueryFindOne, helper.SelectCustom(selectColumn), TableNameUsers, wq), vals...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return &dest, nil
}

func (u *userRepository) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.User, error) {
	var (
		dest []entity.User
	)

	ctx, span := tracer.NewSpan(ctx, "UserRepo.Finds", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	err = u.db.Query(ctx, &dest, fmt.Sprintf(DefaultQueryFinds, helper.SelectCustom(selectColumns), TableNameUsers, wq), vals...)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return dest, nil
}

func (u userRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return u.db.BeginTx(ctx, opts)
}
