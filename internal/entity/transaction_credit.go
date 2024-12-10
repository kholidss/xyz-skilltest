package entity

type TransactionCredit struct {
	ID                 string `json:"id,omitempty" db:"id,omitempty"`
	UserID             string `json:"user_id,omitempty" db:"user_id,omitempty"`
	MerchantID         string `json:"merchant_id,omitempty" db:"merchant_id,omitempty"`
	TransactionID      string `json:"transaction_id,omitempty" db:"transaction_id,omitempty"`
	AssetName          string `json:"asset_name,omitempty" db:"asset_name,omitempty"`
	Tenor              int    `json:"tenor,omitempty" db:"tenor,omitempty"`
	OTRAmount          int    `json:"otr_amount,omitempty" db:"otr_amount,omitempty"`
	FeeAmount          int    `json:"fee_amount,omitempty" db:"fee_amount,omitempty"`
	InterestAmount     int    `json:"interest_amount,omitempty" db:"interest_amount,omitempty"`
	InterestPercentage int    `json:"interest_percentage,omitempty" db:"interest_percentage,omitempty"`
	InstallmentAmount  int    `json:"installment_amount,omitempty" db:"installment_amount,omitempty"`
	DefaultCompleteTimestamp
}
