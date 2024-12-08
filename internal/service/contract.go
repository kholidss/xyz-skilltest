package service

import (
	"context"

	"github.com/kholidss/xyz-skilltest/internal/appctx"
	"github.com/kholidss/xyz-skilltest/internal/entity"
)

type UserService interface {
	ListUser(ctx context.Context) (*[]entity.User, error)
	StoreUser(ctx context.Context) appctx.Response
}
