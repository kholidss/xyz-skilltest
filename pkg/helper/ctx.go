package helper

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/consts"
	"github.com/kholidss/xyz-skilltest/internal/presentation"
)

func GetRequestIDFromFiberCtx(fCtx *fiber.Ctx) string {
	val, ok := fCtx.Locals("request-id").(string)
	if !ok {
		return ""
	}
	return val
}

func GetUserAuthDataFromFiberCtx(fCtx *fiber.Ctx) presentation.UserAuthData {
	val, ok := fCtx.Locals(consts.CtxKeyUserAuthData).(presentation.UserAuthData)
	if !ok {
		return presentation.UserAuthData{}
	}
	return val
}

func SetRequestIDToCtx(ctx context.Context, reqID string) context.Context {
	ctxNew := context.WithValue(ctx, consts.CtxKeyRequestID, reqID)
	return ctxNew
}

func GetRequestIDFromCtx(ctx context.Context) string {
	val, ok := ctx.Value(consts.CtxKeyRequestID).(string)
	if !ok {
		return ""
	}
	return val
}
