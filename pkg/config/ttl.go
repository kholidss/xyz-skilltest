package config

// TTLConfig holds the TTLConfig configuration
type TTLConfig struct {
	TTLTransactionLockCreditUserInMillisecond int `mapstructure:"ttl_transaction_lock_credit_user_in_millisecond"`
}
