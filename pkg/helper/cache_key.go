package helper

import "fmt"

// CacheKeyLockTrxCreditUser is a key of redis cache
func CacheKeyLockTrxCreditUser(userID string) string {
	return fmt.Sprintf("xyz-skilltest:trx-credit-user:%s", userID)
}
