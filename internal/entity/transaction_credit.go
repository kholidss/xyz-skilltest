package entity

type TransactionCredit struct {
	ID                 string `json:"id,omitempty" db:"id,omitempty"`
	UserID             string `json:"user_id,omitempty" db:"user_id,omitempty"`
	MerchantID         string `json:"merchant_id,omitempty" db:"merchant_id,omitempty"`
	TransactionID      string `json:"transaction_id,omitempty" db:"transaction_id,omitempty"`
	AssetName          string `json:"asset_name,omitempty" db:"asset_name,omitempty"`
	LimitAmount        string `json:"limit_amount,omitempty" db:"limit_amount,omitempty"`
	OTRAmount          string `json:"otr_amount,omitempty" db:"otr_amount,omitempty"`
	FeeAmount          string `json:"fee_amount,omitempty" db:"fee_amount,omitempty"`
	InstallmentAmount  string `json:"installment_amount,omitempty" db:"installment_amount,omitempty"`
	InterestAmount     string `json:"interest_amount,omitempty" db:"interest_amount,omitempty"`
	InterestPercentage string `json:"interest_percentage,omitempty" db:"interest_percentage,omitempty"`
	DefaultCompleteTimestamp
}
