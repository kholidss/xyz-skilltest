package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/appctx"
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
	"github.com/kholidss/xyz-skilltest/internal/service"
)

type getAllUser struct {
	service service.UserService
}

func (g *getAllUser) Serve(xCtx appctx.Data) appctx.Response {
	ctx := xCtx.FiberCtx.Context()
	users, err := g.service.ListUser(ctx)
	if err != nil {
		return *appctx.NewResponse().WithError([]appctx.ErrorResp{
			{
				Key:      "PROVIDER_ERR",
				Messages: []string{err.Error()},
			},
		}).WithMessage(err.Error()).WithCode(fiber.StatusBadRequest)
	}

	return *appctx.NewResponse().WithCode(fiber.StatusOK).WithMessage("Success").WithData(users)
}

func NewGetAllUser(svc service.UserService) contract.Controller {
	return &getAllUser{service: svc}
}
