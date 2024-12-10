package transaction

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/kholidss/xyz-skilltest/internal/consts"
	"github.com/kholidss/xyz-skilltest/internal/repositories"
	redislock "github.com/kholidss/xyz-skilltest/pkg/redis_lock"
	"github.com/kholidss/xyz-skilltest/pkg/util"
	"math"
	"strings"
	"time"

	"fmt"
	"github.com/kholidss/xyz-skilltest/internal/presentation"
	"github.com/kholidss/xyz-skilltest/pkg/helper"
	"net/http"

	"github.com/kholidss/xyz-skilltest/internal/appctx"
	"github.com/kholidss/xyz-skilltest/internal/entity"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	"github.com/kholidss/xyz-skilltest/pkg/tracer"
)

func (t *transactionService) CreditUserProcess(ctx context.Context, authData presentation.UserAuthData, payload presentation.ReqTransactionCreditUser) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName("ServiceTransactionCreditUserProcess"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
		)
	)

	ctx, span := tracer.NewSpan(ctx, "service.transaction.credit_user_process", nil)
	defer span.End()

	lf.Append(logger.Any("payload.asset_name", payload.AssetName))
	lf.Append(logger.Any("payload.merchant_id", payload.MerchantID))
	lf.Append(logger.Any("payload.tenor", payload.Tenor))
	lf.Append(logger.Any("payload.otr_amount", payload.OTRAmount))

	//Find merchant
	merchant, err := t.repoMerchant.FindOne(ctx, entity.Merchant{
		ID: payload.MerchantID,
	}, []string{"id", "name"})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("find exist merchant by merchant id got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	if merchant == nil {
		logger.WarnWithContext(ctx, "merchant data not found", lf...)
		return *appctx.NewResponse().WithCode(http.StatusNotFound).WithMessage("Merchant data not found")
	}

	lf.Append(logger.Any("merchant.id", merchant.ID))
	lf.Append(logger.Any("merchant.name", merchant.Name))

	//Find limit by user tenor
	limitCredit, err := t.repoLimit.FindOne(ctx, entity.Limit{
		UserID: authData.UserID,
		Tenor:  payload.Tenor,
	}, []string{"id", "tenor", "limit_amount"})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("find current limit user by tenor got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}
	if limitCredit == nil {
		logger.WarnWithContext(ctx, "current user limit by tenor not exist", lf...)
		return *appctx.NewResponse().WithCode(http.StatusNotFound).WithMessage("Tenor on current user does not exist")
	}

	lf.Append(logger.Any("limit_credit.tenor", limitCredit.Tenor))
	lf.Append(logger.Any("limit_credit.limit_amount", limitCredit.LimitAmount))

	//Do lock current transaction using redislock to avoid concurrent transaction ================================
	err = t.redisLock.Obtain(
		ctx,
		helper.CacheKeyLockTrxCreditUser(authData.UserID),
		t.cfg.TTLConfig.TTLTransactionLockCreditUserInMillisecond,
	)
	if errors.Is(err, redislock.ErrNotObtained) {
		logger.WarnWithContext(ctx, "current credit transaction flagged as concurrent transaction ", lf...)
		return *appctx.NewResponse().WithCode(http.StatusConflict).
			WithMessage("The transaction is flagged as a concurrent transaction. Please retry the operation or ensure no duplicate transactions are being processed")
	}
	defer t.redisLock.ReleaseLock(ctx)
	//Do lock current transaction using redislock to avoid concurrent transaction ================================

	//Calculate final result of installment amount per month
	interestAmount, finalInstallmentAmount := t.calculateInstallmentAmount(payload.OTRAmount, payload.Tenor)

	lf.Append(logger.Any("amount.installment_amount", finalInstallmentAmount))
	lf.Append(logger.Any("amount.admin_fee_amount", t.cfg.CreditConfig.CreditFeeAmount))
	lf.Append(logger.Any("amount.interest_percentage", t.cfg.CreditConfig.CreditInterestPercentage))

	isEligibleToProcess := func() bool {
		return limitCredit.LimitAmount < finalInstallmentAmount
	}()
	if isEligibleToProcess {
		logger.WarnWithContext(ctx, "insufficient loan limit user", lf...)
		return *appctx.NewResponse().WithCode(http.StatusUnprocessableEntity).
			WithMessage(
				fmt.Sprintf(
					"Loan limit is insufficient. The calculated installment is %v, but the remaining loan limit maximum up to %v",
					util.ToFormatRupiah(finalInstallmentAmount),
					util.ToFormatRupiah(limitCredit.LimitAmount),
				),
			)
	}

	// start db transaction
	tx, err := t.repoTransaction.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("start db transaction got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}
	txOpt := repositories.WithTransaction(tx)

	var (
		errTrx         error
		contractNumber = t.generateContractNumber(authData.UserID)
		transactionID  = uuid.New().String()

		remainingLimitAmount = limitCredit.LimitAmount - finalInstallmentAmount
	)

	//Always rollback db transaction if got error process
	defer func() {
		if errTrx != nil && tx != nil {
			_ = tx.Rollback()
		}
	}()

	lf.Append(logger.Any("amount.remaining_limit_amount", remainingLimitAmount))
	lf.Append(logger.Any("transaction.id", transactionID))
	lf.Append(logger.Any("transaction.contract_number", contractNumber))

	//Store transaction data
	errTrx = t.repoTransaction.Store(ctx, entity.Transaction{
		ID:             transactionID,
		UserID:         authData.UserID,
		MerchantID:     payload.MerchantID,
		ContractNumber: contractNumber,
		Type:           consts.TrxTypeCredit,
		TotalAmount:    finalInstallmentAmount,
	}, txOpt)
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("store transaction data got error: %v", errTrx), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	//Store transaction credit data
	errTrx = t.repoTransactionCredit.Store(ctx, entity.TransactionCredit{
		ID:                 uuid.New().String(),
		UserID:             authData.UserID,
		MerchantID:         payload.MerchantID,
		TransactionID:      transactionID,
		AssetName:          payload.AssetName,
		Tenor:              payload.Tenor,
		OTRAmount:          payload.OTRAmount,
		FeeAmount:          t.cfg.CreditConfig.CreditFeeAmount,
		InterestAmount:     interestAmount,
		InterestPercentage: t.cfg.CreditConfig.CreditInterestPercentage,
		InstallmentAmount:  finalInstallmentAmount,
	}, txOpt)
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("store transaction credit data got error: %v", errTrx), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	//Update current limit after calculate the final installment amount
	errTrx = t.repoLimit.Update(ctx, entity.Limit{
		LimitAmount: remainingLimitAmount,
	}, entity.Limit{
		ID: limitCredit.ID,
	}, txOpt)
	if errTrx != nil {
		tracer.AddSpanError(span, errTrx)
		logger.ErrorWithContext(ctx, fmt.Sprintf("update limit amount data got error: %v", errTrx), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	//Commit the db transaction
	if tx != nil {
		_ = tx.Commit()

	}

	logger.InfoWithContext(ctx, "success transaction credit user", lf...)
	return *appctx.NewResponse().
		WithCode(http.StatusCreated).
		WithMessage("Success transaction credit user").
		WithData(presentation.RespTransactionCreditUser{
			Merchant: struct {
				ID   string `json:"id,omitempty"`
				Name string `json:"name,omitempty"`
			}{
				ID:   payload.MerchantID,
				Name: merchant.Name,
			},
			AssetName:            payload.AssetName,
			Tenor:                payload.Tenor,
			OTRAmount:            payload.OTRAmount,
			FeeAmount:            t.cfg.CreditConfig.CreditFeeAmount,
			InterestAmount:       interestAmount,
			InterestPercentage:   t.cfg.CreditConfig.CreditInterestPercentage,
			InstallmentAmount:    finalInstallmentAmount,
			RemainingLimitAmount: remainingLimitAmount,
		})
}

func (t *transactionService) calculateInstallmentAmount(otrAmount int, tenor int) (int, int) {
	otrAmountPerMonth := otrAmount / tenor
	interestAmount := util.CalculatePercentage(otrAmountPerMonth, t.cfg.CreditConfig.CreditInterestPercentage)

	installmentAmount := (otrAmountPerMonth + t.cfg.CreditConfig.CreditFeeAmount) + interestAmount

	return interestAmount, int(math.Ceil(float64(installmentAmount)))
}

func (t *transactionService) generateContractNumber(userID string) string {
	currentTime := time.Now().Format("20060102")
	randomString := util.HashToSha256(util.GenerateRandomString(10) + userID)

	return fmt.Sprintf("%s/%s/%s", consts.TrxTypeCredit, currentTime, strings.ToUpper(randomString[:10]))

}
