package entity

type Merchant struct {
	ID   string `json:"id,omitempty" db:"id,omitempty"`
	Name string `json:"name,omitempty" db:"name,omitempty"`
	Slug string `json:"slug,omitempty" db:"slug,omitempty"`
	DefaultCompleteTimestamp
}
