package config

// CreditConfig holds the CreditConfig configuration
type CreditConfig struct {
	CreditFeeAmount          int `mapstructure:"credit_fee_amount"`
	CreditInterestPercentage int `mapstructure:"credit_interest_percentage"`
}
