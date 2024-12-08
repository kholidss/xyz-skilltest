package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/appctx"
	"github.com/kholidss/xyz-skilltest/internal/bootstrap"
	"github.com/kholidss/xyz-skilltest/internal/controller"
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
	"github.com/kholidss/xyz-skilltest/internal/controller/todo"
	"github.com/kholidss/xyz-skilltest/internal/controller/user"
	"github.com/kholidss/xyz-skilltest/internal/handler"
	"github.com/kholidss/xyz-skilltest/internal/middleware"
	"github.com/kholidss/xyz-skilltest/internal/provider"
	"github.com/kholidss/xyz-skilltest/internal/repositories"
	"github.com/kholidss/xyz-skilltest/internal/service"
	"github.com/kholidss/xyz-skilltest/pkg/config"
)

type router struct {
	cfg   *config.Config
	fiber fiber.Router
}

func (rtr *router) handle(hfn httpHandlerFunc, svc contract.Controller, mdws ...middleware.MiddlewareFunc) fiber.Handler {
	return func(xCtx *fiber.Ctx) error {

		//check registered middleware functions
		if rm := middleware.FilterFunc(rtr.cfg, xCtx, mdws); rm.Code != fiber.StatusOK {
			// return response base on middleware
			res := *appctx.NewResponse().
				WithCode(rm.Code).
				WithError(rm.Errors).
				WithMessage(rm.Message)
			return rtr.response(xCtx, res)
		}

		//send to controller
		resp := hfn(xCtx, svc, rtr.cfg)
		return rtr.response(xCtx, resp)
	}
}

func (rtr *router) response(fiberCtx *fiber.Ctx, resp appctx.Response) error {
	fiberCtx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return fiberCtx.Status(resp.Code).Send(resp.Byte())
}

func (rtr *router) Route() {
	//init db
	db := bootstrap.RegistryMySQLDatabase(rtr.cfg)

	//define repositories
	userRepo := repositories.NewUserRepositoryImpl(db)

	//define services
	userSvc := service.NewUserServiceImpl(userRepo)

	//define middleware
	basicMiddleware := middleware.NewAuthMiddleware()

	//define provider
	example := provider.NewExampleProvider(rtr.cfg)

	//define storage
	//fs := bootstrap.RegistryGCS(rtr.cfg.GCSConfig.ServiceAccountPath)
	//fs := bootstrap.RegistryAWSSession(rtr.cfg)
	//fs := bootstrap.RegistryMinio(rtr.cfg)

	//define cdn
	//cdnStorage := bootstrap.RegistryCDN(rtr.cfg)

	//define controller
	getAllUser := user.NewGetAllUser(userSvc)
	storeUser := user.NewStoreUser(userSvc)
	getTodos := todo.NewGetTodo(example)

	health := controller.NewGetHealth()
	internalV1 := rtr.fiber.Group("/api/internal/v1")

	rtr.fiber.Get("/ping", rtr.handle(
		handler.HttpRequest,
		health,
	))

	internalV1.Get("/users", rtr.handle(
		handler.HttpRequest,
		getAllUser,
		//middleware
		basicMiddleware.Authenticate,
	))

	internalV1.Post("/users", rtr.handle(
		handler.HttpRequest,
		storeUser,
		//middleware
		// basicMiddleware.Authenticate,
	))

	internalV1.Get("/todos", rtr.handle(
		handler.HttpRequest,
		getTodos,
		//middleware
		// basicMiddleware.Authenticate,
	))

}

func NewRouter(cfg *config.Config, fiber fiber.Router) Router {
	return &router{
		cfg:   cfg,
		fiber: fiber,
	}
}
