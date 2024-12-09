package entity

import "time"

type DefaultCompleteTimestamp struct {
	CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	IsDeleted bool       `json:"is_deleted,omitempty" db:"is_deleted,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at,omitempty"`
}
