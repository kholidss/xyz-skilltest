package transaction

import (
	"context"
	"github.com/kholidss/xyz-skilltest/internal/presentation"

	"github.com/kholidss/xyz-skilltest/internal/appctx"
)

type TransactionService interface {
	CreditUserProcess(ctx context.Context, authData presentation.UserAuthData, payload presentation.ReqTransactionCreditUser) appctx.Response
}
