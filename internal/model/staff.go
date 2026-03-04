package model

import (
    "time"
)

type Staff struct {
    ID           string    `gorm:"type:uuid;primaryKey"`
    Username     string    `gorm:"uniqueIndex;not null"`
    PasswordHash string    `gorm:"not null"`
    HospitalID   string    `gorm:"type:uuid;not null"`
    Hospital     Hospital  `gorm:"foreignKey:HospitalID"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}


func (Staff) TableName() string {
    return "staff"
}