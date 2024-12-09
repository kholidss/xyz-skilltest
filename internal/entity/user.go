package entity

type User struct {
	ID           string `json:"id,omitempty" db:"id,omitempty"`
	NIK          string `json:"nik,omitempty" db:"nik,omitempty"`
	FullName     string `json:"full_name,omitempty" db:"full_name,omitempty"`
	LegalName    string `json:"legal_name,omitempty" db:"legal_name,omitempty"`
	PlaceOfBirth string `json:"place_of_birth,omitempty" db:"place_of_birth,omitempty"`
	DateOfBirth  string `json:"date_of_birth,omitempty" db:"date_of_birth,omitempty"`
	Salary       int    `json:"salary,omitempty" db:"salary,omitempty"`
	Password     string `json:"password,omitempty" db:"password,omitempty"`
	DefaultCompleteTimestamp
}
