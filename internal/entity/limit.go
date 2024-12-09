package entity

type Limit struct {
	ID          string `json:"id,omitempty" db:"id,omitempty"`
	UserID      string `json:"user_id,omitempty" db:"user_id,omitempty"`
	Tenor       int    `json:"tenor,omitempty" db:"tenor,omitempty"`
	LimitAmount int    `json:"limit_amount,omitempty" db:"limit_amount,omitempty"`
	DefaultCompleteTimestamp
}
