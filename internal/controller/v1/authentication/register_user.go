package authentication

import (
	"errors"
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
	"github.com/kholidss/xyz-skilltest/pkg/util"
	"github.com/kholidss/xyz-skilltest/pkg/validator"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/appctx"
)

type registerUser struct {
	svcAuthentication authentication.AuthenticationService
}

func NewRegisterUser(svcAuthentication authentication.AuthenticationService) contract.Controller {
	return &registerUser{
		svcAuthentication: svcAuthentication,
	}
}

func (r *registerUser) Serve(xCtx appctx.Data) appctx.Response {
	var (
		requestID = helper.GetRequestIDFromFiberCtx(xCtx.FiberCtx)
		lf        = logger.NewFields(
			logger.EventName("AuthV1RegisterUser"),
			logger.Any("X-Request-ID", requestID),
		)
	)

	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "controller.auth.register_user_v1", nil)
	defer span.End()

	//Inject RequestID to Context
	ctx = helper.SetRequestIDToCtx(ctx, requestID)

	//Parsing the form data request body
	payload, err := r.parse(xCtx.FiberCtx)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("parse payload got error: %v", err), lf...)
		return *appctx.NewResponse().WithMessage(consts.MsgAPIBadRequest).WithCode(fiber.StatusBadRequest)
	}

	lf.Append(logger.Any("payload.nik", masker.Censored(payload.NIK, "*")))
	lf.Append(logger.Any("payload.full_name", payload.FullName))
	lf.Append(logger.Any("payload.legal_name", payload.LegalName))
	lf.Append(logger.Any("payload.place_of_birth", payload.POB))
	lf.Append(logger.Any("payload.dob", payload.DOB))
	lf.Append(logger.Any("payload.salary", payload.Salary))
	lf.Append(logger.Any("payload.password", masker.Censored(payload.Password, "*")))

	// Validate payload
	err = r.validate(payload)
	if err != nil {
		logger.WarnWithContext(ctx, "payload got error validation", lf...)
		return *appctx.NewResponse().WithError(helper.FormatError(err)).
			WithMessage(consts.MsgAPIValidationsError).
			WithCode(fiber.StatusUnprocessableEntity)
	}

	// Validate file
	errFile := r.validateFile(payload)
	if len(errFile) > 0 {
		logger.WarnWithContext(ctx, "payload got error file validation", lf...)
		return *appctx.NewResponse().WithError(errFile).
			WithMessage(consts.MsgAPIValidationsError).
			WithCode(fiber.StatusUnprocessableEntity)
	}

	rsp := r.svcAuthentication.RegisterUser(ctx, payload)
	return rsp
}

func (r *registerUser) parse(c *fiber.Ctx) (presentation.ReqRegisterUser, error) {
	var (
		req presentation.ReqRegisterUser
		err error
	)

	req.NIK = c.FormValue("nik")
	req.FullName = c.FormValue("full_name")
	req.LegalName = c.FormValue("legal_name")
	req.POB = c.FormValue("place_of_birth")
	req.DOB = util.SanitizeDate(c.FormValue("date_of_birth"), "2006-01-02")

	req.Salary, err = strconv.Atoi(func() string {
		if c.FormValue("salary") == "" {
			return "0"
		}
		return c.FormValue("salary")
	}())
	if err != nil {
		return req, errors.New("salary must be a number")
	}

	req.Password = c.FormValue("password")

	//KTP File
	req.FileKTP, err = util.FiberParseFile(c, "ktp")
	if err != nil {
		return req, err
	}

	//Selfie File
	req.FileSelfie, err = util.FiberParseFile(c, "selfie")
	if err != nil {
		return req, err
	}

	return req, nil
}

func (r *registerUser) validate(payload presentation.ReqRegisterUser) error {
	rules := []*validation.FieldRules{
		// NIK
		validation.Field(&payload.NIK, validation.Required, validation.Length(16, 16), validator.ValidateNIK()),

		// FullName
		validation.Field(&payload.FullName, validation.Required, validation.Length(1, 255), validator.ValidateHumanName()),

		// LegalName
		validation.Field(&payload.LegalName, validation.Required, validation.Length(1, 255), validator.ValidateHumanName()),

		// POB
		validation.Field(&payload.POB, validation.Required, validation.Length(1, 100)),

		// DOB
		validation.Field(&payload.DOB, validation.Required, validation.Length(10, 10), validator.ValidateDOB()),

		// Salary
		validation.Field(&payload.Salary, validation.Required),

		// Password
		validation.Field(&payload.Password, validation.Required, validation.Length(8, 50)),

		// FileKTP
		validation.Field(&payload.FileKTP, validation.Required),

		// FileSelfie
		validation.Field(&payload.FileSelfie, validation.Required),
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

func (r *registerUser) validateFile(payload presentation.ReqRegisterUser) []appctx.ErrorResp {
	var (
		errFormat []appctx.ErrorResp

		msgKTP, msgSelfie []string
	)

	if !util.InArray(payload.FileKTP.Mimetype, consts.AllowedMimeTypesKTP) {
		msgKTP = append(
			msgKTP,
			fmt.Sprintf("invalid mimetype ktp, only allow: %s", strings.Join(consts.AllowedMimeTypesKTP, ",")),
		)
	}
	if !util.InArray(payload.FileSelfie.Mimetype, consts.AllowedMimeTypesSelfie) {
		msgSelfie = append(
			msgSelfie,
			fmt.Sprintf("invalid mimetype selfie, only allow: %s", strings.Join(consts.AllowedMimeTypesSelfie, ",")),
		)
	}

	if payload.FileKTP.Size > consts.MaxSizeKTP {
		msgKTP = append(msgKTP,
			fmt.Sprintf(
				"ktp file size %s is too large, max file size is: %s",
				util.HumanFileSize(float64(payload.FileKTP.Size)),
				util.HumanFileSize(float64(consts.MaxSizeKTP))),
		)
	}

	if payload.FileSelfie.Size > consts.MaxSizeSelfie {
		msgSelfie = append(msgSelfie,
			fmt.Sprintf(
				"selfie file size %s is too large, max file size is: %s",
				util.HumanFileSize(float64(payload.FileSelfie.Size)),
				util.HumanFileSize(float64(consts.MaxSizeSelfie))),
		)
	}

	if len(msgKTP) > 0 {
		errFormat = append(errFormat, appctx.ErrorResp{
			Key:      "ktp",
			Messages: msgKTP,
		})
	}

	if len(msgSelfie) > 0 {
		errFormat = append(errFormat, appctx.ErrorResp{
			Key:      "selfie",
			Messages: msgSelfie,
		})
	}

	return errFormat
}
