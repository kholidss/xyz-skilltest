package transaction

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/kholidss/xyz-skilltest/internal/consts"
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
	"github.com/kholidss/xyz-skilltest/internal/presentation"
	"github.com/kholidss/xyz-skilltest/internal/service/transaction"
	"github.com/kholidss/xyz-skilltest/pkg/helper"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	"github.com/kholidss/xyz-skilltest/pkg/tracer"

	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/appctx"
)

type creditUser struct {
	svcTransaction transaction.TransactionService
}

func NewCreditUser(svcTransaction transaction.TransactionService) contract.Controller {
	return &creditUser{
		svcTransaction: svcTransaction,
	}
}

func (c *creditUser) Serve(xCtx appctx.Data) appctx.Response {
	var (
		authInfo  = helper.GetUserAuthDataFromFiberCtx(xCtx.FiberCtx)
		requestID = helper.GetRequestIDFromFiberCtx(xCtx.FiberCtx)
		lf        = logger.NewFields(
			logger.EventName("TransactionV1CreditUser"),
			logger.Any("X-Request-ID", requestID),
		)
		payload presentation.ReqTransactionCreditUser
	)

	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "controller.transaction.credit_user_v1", nil)
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

	lf.Append(logger.Any("payload.asset_name", payload.AssetName))
	lf.Append(logger.Any("payload.merchant_id", payload.MerchantID))
	lf.Append(logger.Any("payload.tenor", payload.Tenor))
	lf.Append(logger.Any("payload.otr_amount", payload.OTRAmount))

	// Validate payload
	err = c.validate(payload)
	if err != nil {
		logger.WarnWithContext(ctx, "payload got error validation", lf...)
		return *appctx.NewResponse().WithError(helper.FormatError(err)).
			WithMessage(consts.MsgAPIValidationsError).
			WithCode(fiber.StatusUnprocessableEntity)
	}

	rsp := c.svcTransaction.CreditUserProcess(ctx, authInfo, payload)
	return rsp
}

func (c *creditUser) validate(payload presentation.ReqTransactionCreditUser) error {
	rules := []*validation.FieldRules{
		// MerchantID
		validation.Field(&payload.MerchantID, validation.Required, is.UUID),

		// AssetName
		validation.Field(&payload.AssetName, validation.Required, validation.Length(1, 255)),

		// Tenor
		validation.Field(&payload.Tenor, validation.Required),

		// OTRAmount
		validation.Field(&payload.OTRAmount, validation.Required),
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
