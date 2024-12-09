package authentication

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kholidss/xyz-skilltest/internal/consts"
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
	"github.com/kholidss/xyz-skilltest/internal/presentation"
	"github.com/kholidss/xyz-skilltest/internal/service/authentication"
	"github.com/kholidss/xyz-skilltest/pkg/helper"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	"github.com/kholidss/xyz-skilltest/pkg/masker"
	"github.com/kholidss/xyz-skilltest/pkg/tracer"
	"github.com/kholidss/xyz-skilltest/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/appctx"
)

type loginUser struct {
	svcAuthentication authentication.AuthenticationService
}

func NewLoginUser(svcAuthentication authentication.AuthenticationService) contract.Controller {
	return &loginUser{
		svcAuthentication: svcAuthentication,
	}
}

func (l *loginUser) Serve(xCtx appctx.Data) appctx.Response {
	var (
		requestID = helper.GetRequestIDFromFiberCtx(xCtx.FiberCtx)
		lf        = logger.NewFields(
			logger.EventName("AuthV1RegisterUser"),
			logger.Any("X-Request-ID", requestID),
		)
		payload presentation.ReqLoginUser
	)

	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "controller.auth.register_user_v1", nil)
	defer span.End()

	//Inject RequestID to Context
	ctx = helper.SetRequestIDToCtx(ctx, requestID)

	//Parsing the form data request body
	err := xCtx.FiberCtx.BodyParser(&payload)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("parse json payload got error: %v", err), lf...)
		return *appctx.NewResponse().WithMessage(consts.MsgAPIBadRequest).WithCode(fiber.StatusBadRequest)
	}

	lf.Append(logger.Any("payload.nik", masker.Censored(payload.NIK, "*")))

	// Validate payload
	err = l.validate(payload)
	if err != nil {
		logger.WarnWithContext(ctx, "payload got error validation", lf...)
		return *appctx.NewResponse().WithError(helper.FormatError(err)).
			WithMessage(consts.MsgAPIValidationsError).
			WithCode(fiber.StatusUnprocessableEntity)
	}

	rsp := l.svcAuthentication.LoginUser(ctx, payload)
	return rsp
}

func (l *loginUser) validate(payload presentation.ReqLoginUser) error {
	rules := []*validation.FieldRules{
		// NIK
		validation.Field(&payload.NIK, validation.Required, validation.Length(16, 16), validator.ValidateNIK()),

		// Password
		validation.Field(&payload.Password, validation.Required),
	}

	err := validation.ValidateStruct(&payload, rules...)
	ve, ok := err.(validation.Errors)
	if !ok {
		ve = make(validation.Errors)
	}

	if len(ve) == 0 {
		return nil
	}

	return ve
}
