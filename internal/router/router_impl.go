package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/appctx"
	"github.com/kholidss/xyz-skilltest/internal/bootstrap"
	"github.com/kholidss/xyz-skilltest/internal/controller"
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
	modulAuthentication "github.com/kholidss/xyz-skilltest/internal/controller/v1/authentication"
	"github.com/kholidss/xyz-skilltest/internal/handler"
	"github.com/kholidss/xyz-skilltest/internal/middleware"
	"github.com/kholidss/xyz-skilltest/internal/repositories"
	"github.com/kholidss/xyz-skilltest/internal/service/authentication"
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
	repoUser := repositories.NewUserRepository(db)
	repoBucket := repositories.NewBucketRepository(db)
	repoLimit := repositories.NewLimitRepository(db)

	//define middleware
	//basicMiddleware := middleware.NewAuthMiddleware()

	//define storage
	//fs := bootstrap.RegistryGCS(rtr.cfg.GCSConfig.ServiceAccountPath)
	//fs := bootstrap.RegistryAWSSession(rtr.cfg)
	//fs := bootstrap.RegistryMinio(rtr.cfg)

	//define cdn
	cdnStorage := bootstrap.RegistryCDN(rtr.cfg)

	//define services
	svcAuthentication := authentication.NewSvcAuthentication(
		rtr.cfg,
		repoUser,
		repoBucket,
		repoLimit,
		cdnStorage,
	)

	//define controller
	ctrlRegisterUser := modulAuthentication.NewRegisterUser(svcAuthentication)
	ctrlLoginUser := modulAuthentication.NewLoginUser(svcAuthentication)

	health := controller.NewGetHealth()
	externalV1 := rtr.fiber.Group("/api/external/v1")

	pathAuthV1 := externalV1.Group("/auth")

	rtr.fiber.Get("/ping", rtr.handle(
		handler.HttpRequest,
		health,
	))

	pathAuthV1.Post("/register-user", rtr.handle(
		handler.HttpRequest,
		ctrlRegisterUser,
	))
	pathAuthV1.Post("/login-user", rtr.handle(
		handler.HttpRequest,
		ctrlLoginUser,
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
