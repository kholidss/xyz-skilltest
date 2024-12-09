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

type transactionCreditRepository struct {
	db mysql.Adapter
}

func NewTransactionCreditRepository(db mysql.Adapter) TransactionCreditRepository {
	return &transactionCreditRepository{
		db: db,
	}
}

func (t transactionCreditRepository) Store(ctx context.Context, payload any, opts ...Option) error {
	var (
		err error
		tx  *sql.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "TransactionCreditRepository", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = t.db.BeginTx(ctx, &sql.TxOptions{
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

	query, val, err := helper.StructQueryInsertMysql(payload, TableNameTransactionCredits, "db", false)

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

func (t transactionCreditRepository) Update(ctx context.Context, payload any, where any, opts ...Option) error {
	var (
		err error
		tx  *sql.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "TransactionCreditRepositorye", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx
	} else {
		tx, err = t.db.BeginTx(ctx, &sql.TxOptions{
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

	q, vals, err := helper.StructToQueryUpdateMysql(payload, where, TableNameTransactionCredits, "db")
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

func (t *transactionCreditRepository) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.TransactionCredit, error) {
	var (
		dest entity.TransactionCredit
	)

	ctx, span := tracer.NewSpan(ctx, "TransactionCreditRepositoryne", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	err = t.db.QueryRow(ctx, &dest, fmt.Sprintf(DefaultQueryFindOne, helper.SelectCustom(selectColumn), TableNameTransactionCredits, wq), vals...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return &dest, nil
}

func (t *transactionCreditRepository) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.TransactionCredit, error) {
	var (
		dest []entity.TransactionCredit
	)

	ctx, span := tracer.NewSpan(ctx, "TransactionCreditRepository", nil)
	defer span.End()

	wq, vals, _, _, err := helper.StructQueryWhereMysql(param, true, "db")
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	err = t.db.Query(ctx, &dest, fmt.Sprintf(DefaultQueryFinds, helper.SelectCustom(selectColumns), TableNameTransactionCredits, wq), vals...)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return dest, nil
}

func (t transactionCreditRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return t.db.BeginTx(ctx, opts)
}
