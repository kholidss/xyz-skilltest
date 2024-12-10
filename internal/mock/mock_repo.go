package mock

import (
	"context"
	"database/sql"
	"github.com/kholidss/xyz-skilltest/internal/entity"
	"github.com/kholidss/xyz-skilltest/internal/repositories"
	"github.com/stretchr/testify/mock"
)

/*
=======================================================================================================
Create Mock Repo User
=======================================================================================================
*/
type MockRepoUser struct {
	mock.Mock
}

func (m *MockRepoUser) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoUser) Update(ctx context.Context, payload any, where any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoUser) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.User, error) {
	args := m.Called(ctx, param, selectColumn)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockRepoUser) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.User, error) {
	args := m.Called(ctx, param, selectColumns)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockRepoUser) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return nil, args.Error(1)
}

/*
=======================================================================================================
Create Mock Repo Merchant
=======================================================================================================
*/
type MockRepoMerchant struct {
	mock.Mock
}

func (m *MockRepoMerchant) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoMerchant) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.Merchant, error) {
	args := m.Called(ctx, param, selectColumn)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Merchant), args.Error(1)
}

func (m *MockRepoMerchant) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Merchant, error) {
	args := m.Called(ctx, param, selectColumns)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.Merchant), args.Error(1)
}

func (m *MockRepoMerchant) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return nil, args.Error(1)
}

/*
=======================================================================================================
Create Mock Repo Bucket
=======================================================================================================
*/
type MockRepoBucket struct {
	mock.Mock
}

func (m *MockRepoBucket) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoBucket) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.Bucket, error) {
	args := m.Called(ctx, param, selectColumn)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Bucket), args.Error(1)
}

func (m *MockRepoBucket) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Bucket, error) {
	args := m.Called(ctx, param, selectColumns)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.Bucket), args.Error(1)
}

func (m *MockRepoBucket) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return nil, args.Error(1)
}

/*
=======================================================================================================
Create Mock Repo Limit
=======================================================================================================
*/

type MockRepoLimit struct {
	mock.Mock
}

func (m *MockRepoLimit) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoLimit) Update(ctx context.Context, payload any, where any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, where, opts)
	return args.Error(0)
}

func (m *MockRepoLimit) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.Limit, error) {
	args := m.Called(ctx, param, selectColumn)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Limit), args.Error(1)
}

func (m *MockRepoLimit) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Limit, error) {
	args := m.Called(ctx, param, selectColumns)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.Limit), args.Error(1)
}

func (m *MockRepoLimit) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return nil, args.Error(1)
}

/*
=======================================================================================================
Create Mock Repo Transaction
=======================================================================================================
*/
type MockRepoTransaction struct {
	mock.Mock
}

func (m *MockRepoTransaction) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoTransaction) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.Transaction, error) {
	args := m.Called(ctx, param, selectColumn)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockRepoTransaction) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Transaction, error) {
	args := m.Called(ctx, param, selectColumns)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func (m *MockRepoTransaction) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return nil, args.Error(1)
}

/*
=======================================================================================================
Create Mock Repo Transaction Credit
=======================================================================================================
*/
type MockRepoTransactionCredit struct {
	mock.Mock
}

func (m *MockRepoTransactionCredit) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoTransactionCredit) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.TransactionCredit, error) {
	args := m.Called(ctx, param, selectColumn)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.TransactionCredit), args.Error(1)
}

func (m *MockRepoTransactionCredit) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.TransactionCredit, error) {
	args := m.Called(ctx, param, selectColumns)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.TransactionCredit), args.Error(1)
}

func (m *MockRepoTransactionCredit) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return nil, args.Error(1)
}
