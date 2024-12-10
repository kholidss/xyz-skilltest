package transaction

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kholidss/xyz-skilltest/internal/entity"
	mockMethod "github.com/kholidss/xyz-skilltest/internal/mock"
	"github.com/kholidss/xyz-skilltest/internal/presentation"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	redislock "github.com/kholidss/xyz-skilltest/pkg/redis_lock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRegisterService(t *testing.T) {
	testCases := []struct {
		testName      string
		inputPayload  presentation.ReqTransactionCreditUser
		inputAuthData presentation.UserAuthData

		rEFindOneMerchant        []any
		rEFindOneLimit           []any
		rEStoreTransaction       []any
		rEStoreTransactionCredit []any
		rEUpdateLimit            []any
		rEBeginTxTransaction     []any

		rERedisLockObtain []any

		expectedHTTPCode int
	}{
		{
			testName: "[TEST] Valid response and success credit transaction",
			inputPayload: presentation.ReqTransactionCreditUser{
				MerchantID: "mid-uuid",
				AssetName:  "Motor Supra",
				Tenor:      6,
				OTRAmount:  15000000,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaajjjj",
				FullName:    "John Doe",
				LegalName:   "John Doe",
			},

			rEFindOneMerchant: []any{&entity.Merchant{
				ID:   "mid-uuid",
				Name: "PT Xyz",
			}, nil},
			rEFindOneLimit: []any{&entity.Limit{
				ID:          uuid.New().String(),
				UserID:      "",
				Tenor:       6,
				LimitAmount: 20000000,
			}, nil},
			rERedisLockObtain:        []any{nil},
			rEBeginTxTransaction:     []any{nil, nil},
			rEStoreTransaction:       []any{nil},
			rEStoreTransactionCredit: []any{nil},
			rEUpdateLimit:            []any{nil},
			expectedHTTPCode:         fiber.StatusCreated,
		},
		{
			testName: "[TEST] Error FindOne merchant",
			inputPayload: presentation.ReqTransactionCreditUser{
				MerchantID: "mid-uuid",
				AssetName:  "Motor Supra",
				Tenor:      6,
				OTRAmount:  15000000,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaajjjj",
				FullName:    "John Doe",
				LegalName:   "John Doe",
			},

			rEFindOneMerchant: []any{nil, errors.New("some error")},
			rEFindOneLimit: []any{&entity.Limit{
				ID:          uuid.New().String(),
				UserID:      "",
				Tenor:       6,
				LimitAmount: 20000000,
			}, nil},
			rERedisLockObtain:        []any{nil},
			rEBeginTxTransaction:     []any{nil, nil},
			rEStoreTransaction:       []any{nil},
			rEStoreTransactionCredit: []any{nil},
			rEUpdateLimit:            []any{nil},
			expectedHTTPCode:         fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Merchant not found",
			inputPayload: presentation.ReqTransactionCreditUser{
				MerchantID: "mid-uuid",
				AssetName:  "Motor Supra",
				Tenor:      6,
				OTRAmount:  15000000,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaajjjj",
				FullName:    "John Doe",
				LegalName:   "John Doe",
			},

			rEFindOneMerchant: []any{nil, nil},
			rEFindOneLimit: []any{&entity.Limit{
				ID:          uuid.New().String(),
				UserID:      "",
				Tenor:       6,
				LimitAmount: 20000000,
			}, nil},
			rERedisLockObtain:        []any{nil},
			rEBeginTxTransaction:     []any{nil, nil},
			rEStoreTransaction:       []any{nil},
			rEStoreTransactionCredit: []any{nil},
			rEUpdateLimit:            []any{nil},
			expectedHTTPCode:         fiber.StatusNotFound,
		},
		{
			testName: "[TEST] Error FindOne limit",
			inputPayload: presentation.ReqTransactionCreditUser{
				MerchantID: "mid-uuid",
				AssetName:  "Motor Supra",
				Tenor:      6,
				OTRAmount:  15000000,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaajjjj",
				FullName:    "John Doe",
				LegalName:   "John Doe",
			},

			rEFindOneMerchant: []any{&entity.Merchant{
				ID:   "mid-uuid",
				Name: "PT Xyz",
			}, nil},
			rEFindOneLimit:           []any{nil, errors.New("some error")},
			rERedisLockObtain:        []any{nil},
			rEBeginTxTransaction:     []any{nil, nil},
			rEStoreTransaction:       []any{nil},
			rEStoreTransactionCredit: []any{nil},
			rEUpdateLimit:            []any{nil},
			expectedHTTPCode:         fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Limit by tenor got not found",
			inputPayload: presentation.ReqTransactionCreditUser{
				MerchantID: "mid-uuid",
				AssetName:  "Motor Supra",
				Tenor:      6,
				OTRAmount:  15000000,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaajjjj",
				FullName:    "John Doe",
				LegalName:   "John Doe",
			},

			rEFindOneMerchant: []any{&entity.Merchant{
				ID:   "mid-uuid",
				Name: "PT Xyz",
			}, nil},
			rEFindOneLimit:           []any{nil, nil},
			rERedisLockObtain:        []any{nil},
			rEBeginTxTransaction:     []any{nil, nil},
			rEStoreTransaction:       []any{nil},
			rEStoreTransactionCredit: []any{nil},
			rEUpdateLimit:            []any{nil},
			expectedHTTPCode:         fiber.StatusNotFound,
		},
		{
			testName: "[TEST] Got concurrent user transaction",
			inputPayload: presentation.ReqTransactionCreditUser{
				MerchantID: "mid-uuid",
				AssetName:  "Motor Supra",
				Tenor:      6,
				OTRAmount:  15000000,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaajjjj",
				FullName:    "John Doe",
				LegalName:   "John Doe",
			},

			rEFindOneMerchant: []any{&entity.Merchant{
				ID:   "mid-uuid",
				Name: "PT Xyz",
			}, nil},
			rEFindOneLimit: []any{&entity.Limit{
				ID:          uuid.New().String(),
				UserID:      "",
				Tenor:       6,
				LimitAmount: 20000000,
			}, nil},
			rERedisLockObtain:        []any{redislock.ErrNotObtained},
			rEBeginTxTransaction:     []any{nil, nil},
			rEStoreTransaction:       []any{nil},
			rEStoreTransactionCredit: []any{nil},
			rEUpdateLimit:            []any{nil},
			expectedHTTPCode:         fiber.StatusConflict,
		},
		{
			testName: "[TEST] Insufficient user credit limit amount",
			inputPayload: presentation.ReqTransactionCreditUser{
				MerchantID: "mid-uuid",
				AssetName:  "Motor Supra",
				Tenor:      6,
				OTRAmount:  15000000,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaajjjj",
				FullName:    "John Doe",
				LegalName:   "John Doe",
			},

			rEFindOneMerchant: []any{&entity.Merchant{
				ID:   "mid-uuid",
				Name: "PT Xyz",
			}, nil},
			rEFindOneLimit: []any{&entity.Limit{
				ID:          uuid.New().String(),
				UserID:      "",
				Tenor:       6,
				LimitAmount: 200000,
			}, nil},
			rERedisLockObtain:        []any{nil},
			rEBeginTxTransaction:     []any{nil, nil},
			rEStoreTransaction:       []any{nil},
			rEStoreTransactionCredit: []any{nil},
			rEUpdateLimit:            []any{nil},
			expectedHTTPCode:         fiber.StatusUnprocessableEntity,
		},
		{
			testName: "[TEST] Error begin db transaction",
			inputPayload: presentation.ReqTransactionCreditUser{
				MerchantID: "mid-uuid",
				AssetName:  "Motor Supra",
				Tenor:      6,
				OTRAmount:  15000000,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaajjjj",
				FullName:    "John Doe",
				LegalName:   "John Doe",
			},

			rEFindOneMerchant: []any{&entity.Merchant{
				ID:   "mid-uuid",
				Name: "PT Xyz",
			}, nil},
			rEFindOneLimit: []any{&entity.Limit{
				ID:          uuid.New().String(),
				UserID:      "",
				Tenor:       6,
				LimitAmount: 20000000,
			}, nil},
			rERedisLockObtain:        []any{nil},
			rEBeginTxTransaction:     []any{nil, errors.New("some error")},
			rEStoreTransaction:       []any{nil},
			rEStoreTransactionCredit: []any{nil},
			rEUpdateLimit:            []any{nil},
			expectedHTTPCode:         fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Error store transaction data",
			inputPayload: presentation.ReqTransactionCreditUser{
				MerchantID: "mid-uuid",
				AssetName:  "Motor Supra",
				Tenor:      6,
				OTRAmount:  15000000,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaajjjj",
				FullName:    "John Doe",
				LegalName:   "John Doe",
			},

			rEFindOneMerchant: []any{&entity.Merchant{
				ID:   "mid-uuid",
				Name: "PT Xyz",
			}, nil},
			rEFindOneLimit: []any{&entity.Limit{
				ID:          uuid.New().String(),
				UserID:      "",
				Tenor:       6,
				LimitAmount: 20000000,
			}, nil},
			rERedisLockObtain:        []any{nil},
			rEBeginTxTransaction:     []any{nil, nil},
			rEStoreTransaction:       []any{errors.New("some error")},
			rEStoreTransactionCredit: []any{nil},
			rEUpdateLimit:            []any{nil},
			expectedHTTPCode:         fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Error store transaction credit data",
			inputPayload: presentation.ReqTransactionCreditUser{
				MerchantID: "mid-uuid",
				AssetName:  "Motor Supra",
				Tenor:      6,
				OTRAmount:  15000000,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaajjjj",
				FullName:    "John Doe",
				LegalName:   "John Doe",
			},

			rEFindOneMerchant: []any{&entity.Merchant{
				ID:   "mid-uuid",
				Name: "PT Xyz",
			}, nil},
			rEFindOneLimit: []any{&entity.Limit{
				ID:          uuid.New().String(),
				UserID:      "",
				Tenor:       6,
				LimitAmount: 20000000,
			}, nil},
			rERedisLockObtain:        []any{nil},
			rEBeginTxTransaction:     []any{nil, nil},
			rEStoreTransaction:       []any{nil},
			rEStoreTransactionCredit: []any{errors.New("some error")},
			rEUpdateLimit:            []any{nil},
			expectedHTTPCode:         fiber.StatusInternalServerError,
		},
		{
			testName: "[TEST] Error update limit amount",
			inputPayload: presentation.ReqTransactionCreditUser{
				MerchantID: "mid-uuid",
				AssetName:  "Motor Supra",
				Tenor:      6,
				OTRAmount:  15000000,
			},
			inputAuthData: presentation.UserAuthData{
				UserID:      uuid.New().String(),
				AccessToken: "xxxaaaajjjj",
				FullName:    "John Doe",
				LegalName:   "John Doe",
			},

			rEFindOneMerchant: []any{&entity.Merchant{
				ID:   "mid-uuid",
				Name: "PT Xyz",
			}, nil},
			rEFindOneLimit: []any{&entity.Limit{
				ID:          uuid.New().String(),
				UserID:      "",
				Tenor:       6,
				LimitAmount: 20000000,
			}, nil},
			rERedisLockObtain:        []any{nil},
			rEBeginTxTransaction:     []any{nil, nil},
			rEStoreTransaction:       []any{nil},
			rEStoreTransactionCredit: []any{nil},
			rEUpdateLimit:            []any{errors.New("some error")},
			expectedHTTPCode:         fiber.StatusInternalServerError,
		},
	}

	cfg, err := config.LoadAllConfigs()
	if err != nil {
		t.Errorf("init config got error: %v", err)
	}

	for _, test := range testCases {
		var (
			mockRepoMerchant          = new(mockMethod.MockRepoMerchant)
			mockRepoLimit             = new(mockMethod.MockRepoLimit)
			mockRepoTransaction       = new(mockMethod.MockRepoTransaction)
			mockRepoTransactionCredit = new(mockMethod.MockRepoTransactionCredit)

			mockRedisLock = new(mockMethod.MockRedisLock)
		)

		// Declare on mock with method ============================================================
		mockRepoMerchant.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(test.rEFindOneMerchant...)

		mockRepoLimit.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(test.rEFindOneLimit...)
		mockRepoLimit.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(test.rEUpdateLimit...)

		mockRepoTransaction.On("BeginTx", mock.Anything, mock.Anything).Return(test.rEBeginTxTransaction...)
		mockRepoTransaction.On("Store", mock.Anything, mock.Anything, mock.Anything).Return(test.rEStoreTransaction...)
		mockRepoTransactionCredit.On("Store", mock.Anything, mock.Anything, mock.Anything).Return(test.rEStoreTransactionCredit...)

		mockRedisLock.On("Obtain", mock.Anything, mock.Anything, mock.Anything).Return(test.rERedisLockObtain...)
		mockRedisLock.On("ReleaseLock", mock.Anything).Return(nil)

		// ========================================================================================

		// Initiate service
		e := NewSvcTransaction(
			cfg,
			mockRepoMerchant,
			mockRepoLimit,
			mockRepoTransaction,
			mockRepoTransactionCredit,
			mockRedisLock,
		)

		resp := e.CreditUserProcess(context.Background(), test.inputAuthData, test.inputPayload)

		assert.Equal(t, test.expectedHTTPCode, resp.Code, fmt.Sprintf(
			`Error test : %s. Expectation http code %d but got %d`,
			test.testName,
			test.expectedHTTPCode,
			resp.Code))

	}

}
