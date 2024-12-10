package transaction

import (
	"github.com/kholidss/xyz-skilltest/internal/repositories"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	redislock "github.com/kholidss/xyz-skilltest/pkg/redis_lock"
)

type transactionService struct {
	cfg                   *config.Config
	repoMerchant          repositories.MerchantRepository
	repoLimit             repositories.LimitRepository
	repoTransaction       repositories.TransactionRepository
	repoTransactionCredit repositories.TransactionCreditRepository
	redisLock             redislock.Locker
}

func NewSvcTransaction(
	cfg *config.Config,
	repoMerchant repositories.MerchantRepository,
	repoLimit repositories.LimitRepository,
	repoTransaction repositories.TransactionRepository,
	repoTransactionCredit repositories.TransactionCreditRepository,
	redisLock redislock.Locker,
) TransactionService {
	return &transactionService{
		cfg:                   cfg,
		repoMerchant:          repoMerchant,
		repoLimit:             repoLimit,
		repoTransaction:       repoTransaction,
		repoTransactionCredit: repoTransactionCredit,
		redisLock:             redisLock,
	}
}
