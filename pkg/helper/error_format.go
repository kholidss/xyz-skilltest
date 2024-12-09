package helper

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kholidss/xyz-skilltest/internal/appctx"
)

func FormatError(err error) []appctx.ErrorResp {
	var errors []appctx.ErrorResp

	if err == nil {
		return nil
	}

	switch ev := err.(type) {
	case validation.Errors:
		return formatOzzoValidationError(ev)
	default:
		return append(errors, appctx.ErrorResp{
			Key:      "-",
			Messages: []string{err.Error()},
		})
	}
}

func formatOzzoValidationError(err validation.Errors) []appctx.ErrorResp {
	var errors []appctx.ErrorResp

	for k, e := range err {
		errors = append(errors, appctx.ErrorResp{
			Key:      k,
			Messages: []string{e.Error()},
		})
	}

	return errors
}
