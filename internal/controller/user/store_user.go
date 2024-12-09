package user

import (
	"github.com/kholidss/xyz-skilltest/internal/appctx"
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
	"github.com/kholidss/xyz-skilltest/internal/service/authentication"
	"github.com/kholidss/xyz-skilltest/pkg/tracer"
)

type storeUser struct {
	service authentication.UserService
}

func (g *storeUser) Serve(xCtx appctx.Data) appctx.Response {
	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "Controller.CreateUser", nil)
	defer span.End()

	res := g.service.StoreUser(ctx)
	return res
}

func NewStoreUser(svc authentication.UserService) contract.Controller {
	return &storeUser{service: svc}
}
