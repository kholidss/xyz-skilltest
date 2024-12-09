package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/appctx"
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
)

type getHealth struct {
}

func (g *getHealth) Serve(xCtx appctx.Data) appctx.Response {
	// Ping Endpoint
	return *appctx.NewResponse().WithCode(fiber.StatusOK).WithMessage("ok").WithData(struct {
		Message string `json:"message"`
	}{
		Message: "OK!",
	})
}

func NewGetHealth() contract.Controller {
	return &getHealth{}
}
