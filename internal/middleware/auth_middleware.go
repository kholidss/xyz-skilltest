package middleware

import (
	"fmt"
	"github.com/kholidss/xyz-skilltest/internal/entity"
	"github.com/kholidss/xyz-skilltest/internal/presentation"
	"github.com/kholidss/xyz-skilltest/internal/repositories"
	"github.com/kholidss/xyz-skilltest/pkg/helper"
	"github.com/spf13/cast"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/appctx"
	"github.com/kholidss/xyz-skilltest/internal/consts"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	"github.com/kholidss/xyz-skilltest/pkg/httpclient"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	"github.com/kholidss/xyz-skilltest/pkg/tracer"
	"github.com/kholidss/xyz-skilltest/pkg/util"
)

type userAuth struct {
	cfg      *config.Config
	repoUser repositories.UserRepository
}

func NewUserAuthMiddleware(cfg *config.Config, repoUser repositories.UserRepository) *userAuth {
	return &userAuth{
		cfg:      cfg,
		repoUser: repoUser,
	}
}

func (u *userAuth) Authenticate(xCtx *fiber.Ctx) appctx.Response {
	var (
		requestID = helper.GetRequestIDFromFiberCtx(xCtx)
		lf        = logger.NewFields(
			logger.EventName("UserAuthMiddleware"),
			logger.Any("X-Request-ID", requestID),
		)
	)

	ctx, span := tracer.NewSpan(xCtx.Context(), "middleware.user_auth", nil)
	defer span.End()

	headerAuth := xCtx.Get(httpclient.Authorization)

	if headerAuth == "" {
		logger.WarnWithContext(ctx, "authorization header is missing", lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage(consts.MsgAPIUnautorized)
	}

	if !strings.HasPrefix(headerAuth, "Bearer ") {
		logger.WarnWithContext(ctx, "Authorization header is missing 'Bearer' prefix", lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage(consts.MsgAPIUnautorized)
	}

	token := strings.TrimPrefix(headerAuth, "Bearer ")

	// Validate the JWT token using the public key
	_, claims, err := util.ValidateJWT(token, u.cfg.AppConfig.APPPublicKey)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.WarnWithContext(ctx, "failed to validate user JWT token", lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage(consts.MsgAPIUnautorized)
	}

	// Parse user claims
	userClaim, ok := claims["user"].(map[string]any)
	if !ok {
		logger.WarnWithContext(ctx, "user claims not found", lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage(consts.MsgAPIUnautorized)
	}

	var (
		userID = cast.ToString(userClaim["user_ids"])
	)

	//Fetch user data
	user, err := u.repoUser.FindOne(ctx, entity.User{
		ID: userID,
	}, []string{"id", "full_name", "legal_name"})

	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("fetch user data got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage(consts.MsgAPIUnautorized)
	}
	if user == nil {
		logger.WarnWithContext(ctx, "got null result of user data", lf...)
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage(consts.MsgAPIUnautorized)
	}

	//Inject user auth data to fiber context
	xCtx.Locals(consts.CtxKeyUserAuthData, presentation.UserAuthData{
		UserID:      cast.ToString(userClaim["user_id"]),
		AccessToken: token,
		FullName:    user.FullName,
		LegalName:   user.LegalName,
	})

	return *appctx.NewResponse().WithCode(fiber.StatusOK)
}
