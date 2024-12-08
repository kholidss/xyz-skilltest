package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/appctx"
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
	"github.com/kholidss/xyz-skilltest/pkg/config"
)

func HttpRequest(xCtx *fiber.Ctx, svc contract.Controller, conf *config.Config) appctx.Response {
	data := appctx.Data{
		FiberCtx: xCtx,
		Cfg:      conf,
	}
	return svc.Serve(data)
}
