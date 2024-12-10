package repositories

import (
	"context"
	"database/sql"
	"github.com/kholidss/xyz-skilltest/internal/entity"
)

type UserRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumn []string) (*entity.User, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.User, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type MerchantRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumn []string) (*entity.Merchant, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Merchant, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type BucketRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumn []string) (*entity.Bucket, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Bucket, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type LimitRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	Update(ctx context.Context, payload any, where any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumn []string) (*entity.Limit, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Limit, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type TransactionRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumn []string) (*entity.Transaction, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Transaction, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type TransactionCreditRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumn []string) (*entity.TransactionCredit, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.TransactionCredit, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
