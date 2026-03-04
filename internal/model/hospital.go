package model

import (
	"time"
)

type Hospital struct {
	ID string `gorm:"type:uuid;primarykey"`
	Name string `gorm:"not null"`
	Code string `gorm:"uniqueIndex;not null"`
	APIBaseURL string `gorm:"column:api_base_url;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}