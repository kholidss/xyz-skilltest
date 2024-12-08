package todo

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/appctx"
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
	"github.com/kholidss/xyz-skilltest/internal/provider"
)

type todo struct {
	provider provider.Example
}

func (t todo) Serve(xCtx appctx.Data) appctx.Response {
	todos, err := t.provider.GetTodos(xCtx.FiberCtx.Context())
	if err != nil {
		if err != nil {
			return *appctx.NewResponse().WithError([]appctx.ErrorResp{
				{
					Key:      "PROVIDER_ERR",
					Messages: []string{err.Error()},
				},
			}).WithMessage("Get Data Failed").WithCode(fiber.StatusBadRequest)
		}
	}

	var result = make([]map[string]interface{}, 0)

	json.Unmarshal(todos, &result)

	return *appctx.NewResponse().WithCode(fiber.StatusOK).WithData(result)
}

func NewGetTodo(provider provider.Example) contract.Controller {
	return &todo{provider: provider}
}
