package authentication

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/kholidss/xyz-skilltest/internal/consts"
	"github.com/kholidss/xyz-skilltest/internal/presentation"
	"github.com/kholidss/xyz-skilltest/internal/repositories"
	"github.com/kholidss/xyz-skilltest/pkg/cipher"
	"github.com/kholidss/xyz-skilltest/pkg/helper"
	"github.com/kholidss/xyz-skilltest/pkg/limit_schema"
	"github.com/kholidss/xyz-skilltest/pkg/masker"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strings"

	"github.com/kholidss/xyz-skilltest/internal/appctx"
	"github.com/kholidss/xyz-skilltest/internal/entity"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	"github.com/kholidss/xyz-skilltest/pkg/tracer"
)

func (a authenticationService) RegisterUser(ctx context.Context, payload presentation.ReqRegisterUser) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName("ServiceAuthRegisterUser"),
			logger.Any("X-Request-ID", helper.GetRequestIDFromCtx(ctx)),
		)
	)

	ctx, span := tracer.NewSpan(ctx, "service.auth.register_user", nil)
	defer span.End()

	lf.Append(logger.Any("payload.nik", masker.Censored(payload.NIK, "*")))
	lf.Append(logger.Any("payload.full_name", payload.FullName))
	lf.Append(logger.Any("payload.legal_name", payload.LegalName))
	lf.Append(logger.Any("payload.place_of_birth", payload.POB))
	lf.Append(logger.Any("payload.dob", payload.DOB))
	lf.Append(logger.Any("payload.salary", payload.Salary))
	lf.Append(logger.Any("payload.password", masker.Censored(payload.Password, "*")))

	//Find exist NIK
	user, err := a.repoUser.FindOne(ctx, entity.User{
		NIK: payload.NIK,
	}, []string{"id", "full_name", "legal_name", "created_at", "updated_at"})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("find exist user by nik got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	if user != nil {
		lf.Append(logger.Any("exist_user.id", user.ID))
		lf.Append(logger.Any("exist_user.full_name", user.FullName))
		lf.Append(logger.Any("exist_user.legal_name", user.LegalName))
		lf.Append(logger.Any("exist_user.created_at", user.CreatedAt))
		lf.Append(logger.Any("exist_user.updated_at", user.UpdatedAt))

		logger.WarnWithContext(ctx, "user by nik already registered", lf...)
		return *appctx.NewResponse().WithCode(http.StatusUnprocessableEntity).WithMessage("NIK already registered")
	}

	var (
		objectKTP, objectSelfie *uploader.UploadResult

		fileNameKTP, pathKTP       = helper.GeneratePathAndFilenameStorage("ktp", strings.Split(payload.FileKTP.Mimetype, "/")[1])
		fileNameSelfie, pathSelfie = helper.GeneratePathAndFilenameStorage("selfie", strings.Split(payload.FileSelfie.Mimetype, "/")[1])
	)

	gr, _ := errgroup.WithContext(ctx)

	//Stream and upload file KTP and Selfie
	gr.Go(func() error {
		rs, err := a.cdn.Put(ctx, pathKTP, payload.FileKTP.File)
		v, ok := rs.(*uploader.UploadResult)
		if ok {
			objectKTP = v
		}
		return err
	})
	gr.Go(func() error {
		rs, err := a.cdn.Put(ctx, pathSelfie, payload.FileSelfie.File)
		v, ok := rs.(*uploader.UploadResult)
		if ok {
			objectSelfie = v
		}
		return err
	})
	err = gr.Wait()
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("upload ktp or selfie to storage got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	// start db transaction
	tx, err := a.repoUser.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("start db transaction got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	txOpt := repositories.WithTransaction(tx)

	passwordUser, err := cipher.EncryptAES256(payload.Password, a.cfg.AppConfig.AppPasswordSecret)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("encrypt user password with aes256 methode got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	var (
		userID   = uuid.New().String()
		errStore error
	)

	//Always rollback db transaction if got error process
	defer func() {
		if errStore != nil && tx != nil {
			_ = tx.Rollback()
		}
	}()

	lf.Append(logger.Any("result.user_id", userID))

	//Store user data
	errStore = a.repoUser.Store(ctx, entity.User{
		ID:           userID,
		NIK:          payload.NIK,
		FullName:     payload.FullName,
		LegalName:    payload.LegalName,
		PlaceOfBirth: payload.POB,
		DateOfBirth:  payload.DOB,
		Salary:       payload.Salary,
		Password:     passwordUser,
	}, txOpt)
	if errStore != nil {
		tracer.AddSpanError(span, errStore)
		logger.ErrorWithContext(ctx, fmt.Sprintf("store user data got error: %v", errStore), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	jsonStorageBuilder := func(obj *uploader.UploadResult) json.RawMessage {
		if obj == nil {
			return nil
		}
		v, err := json.Marshal(obj)
		if err != nil {
			return nil
		}
		return v
	}

	//Store bucket user ktp data
	errStore = a.repoBucket.Store(ctx, entity.Bucket{
		ID:             uuid.New().String(),
		Filename:       fileNameKTP,
		IdentifierID:   userID,
		IdentifierName: consts.BucketIdentifierUserKTP,
		Mimetype:       payload.FileKTP.Mimetype,
		Provider:       strings.ToLower(a.cfg.CDNConfig.Provider),
		URL:            objectKTP.URL,
		Path:           pathKTP,
		SupportData:    jsonStorageBuilder(objectKTP),
	}, txOpt)
	if errStore != nil {
		tracer.AddSpanError(span, errStore)
		logger.ErrorWithContext(ctx, fmt.Sprintf("store bucket user ktp data got error: %v", errStore), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	//Store bucket user selfie data
	errStore = a.repoBucket.Store(ctx, entity.Bucket{
		ID:             uuid.New().String(),
		Filename:       fileNameSelfie,
		IdentifierID:   userID,
		IdentifierName: consts.BucketIdentifierUserSelfie,
		Mimetype:       payload.FileSelfie.Mimetype,
		Provider:       strings.ToLower(a.cfg.CDNConfig.Provider),
		URL:            objectSelfie.URL,
		Path:           pathSelfie,
		SupportData:    jsonStorageBuilder(objectSelfie),
	}, txOpt)
	if errStore != nil {
		tracer.AddSpanError(span, errStore)
		logger.ErrorWithContext(ctx, fmt.Sprintf("store bucket user selfie data got error: %v", errStore), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
	}

	//Store limit schema by salary
	limitSchema := limit_schema.BuildLimitSchema(payload.Salary)
	lf.Append(logger.Any("result.limit_schema", limitSchema))
	for k, v := range limitSchema {
		errStore = a.repoLimit.Store(ctx, entity.Limit{
			ID:          uuid.New().String(),
			UserID:      userID,
			Tenor:       k,
			LimitAmount: v,
		}, txOpt)
		if errStore != nil {
			tracer.AddSpanError(span, errStore)
			logger.ErrorWithContext(ctx, fmt.Sprintf("store limit user data got error: %v", errStore), lf...)
			return *appctx.NewResponse().WithCode(http.StatusInternalServerError)
		}
	}

	//Commit the db transaction
	if tx != nil {
		_ = tx.Commit()
	}

	logger.InfoWithContext(ctx, "success register user", lf...)
	return *appctx.NewResponse().
		WithCode(http.StatusCreated).
		WithMessage("Success register user").
		WithData(presentation.RespRegisterUser{
			UserID:    userID,
			FullName:  payload.FullName,
			LegalName: payload.LegalName,
			DOB:       payload.DOB,
		})

}
