package entity

type Transaction struct {
	ID             string `json:"id,omitempty" db:"id,omitempty"`
	UserID         string `json:"user_id,omitempty" db:"user_id,omitempty"`
	MerchantID     string `json:"merchant_id,omitempty" db:"merchant_id,omitempty"`
	ContractNumber string `json:"contract_number,omitempty" db:"contract_number,omitempty"`
	Type           string `json:"type,omitempty" db:"type,omitempty"`
	TotalAmount    int    `json:"total_amount,omitempty" db:"total_amount,omitempty"`
	DefaultCompleteTimestamp
}
