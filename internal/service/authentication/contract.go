package authentication

import (
	"context"
	"github.com/kholidss/xyz-skilltest/internal/presentation"

	"github.com/kholidss/xyz-skilltest/internal/appctx"
)

type AuthenticationService interface {
	RegisterUser(ctx context.Context, payload presentation.ReqRegisterUser) appctx.Response
}
