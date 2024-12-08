package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/appctx"
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
	"github.com/kholidss/xyz-skilltest/pkg/config"
)

type httpHandlerFunc func(xCtx *fiber.Ctx, svc contract.Controller, conf *config.Config) appctx.Response

type Router interface {
	Route()
}
