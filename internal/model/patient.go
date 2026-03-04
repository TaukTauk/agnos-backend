package model

import (
    "time"
)

type Patient struct {
    ID           string    `gorm:"type:uuid;primaryKey"`
    HospitalID   string    `gorm:"type:uuid;not null;uniqueIndex:idx_hospital_national;uniqueIndex:idx_hospital_passport"`
    Hospital     Hospital  `gorm:"foreignKey:HospitalID"`
    FirstNameTH  string    `gorm:"column:first_name_th"`
    MiddleNameTH string    `gorm:"column:middle_name_th"`
    LastNameTH   string    `gorm:"column:last_name_th"`
    FirstNameEN  string    `gorm:"column:first_name_en"`
    MiddleNameEN string    `gorm:"column:middle_name_en"`
    LastNameEN   string    `gorm:"column:last_name_en"`
    DateOfBirth  string    `gorm:"type:date"`
    PatientHN    string    `gorm:"column:patient_hn"`
    NationalID   *string   `gorm:"uniqueIndex:idx_hospital_national"`
    PassportID   *string   `gorm:"uniqueIndex:idx_hospital_passport"`
    PhoneNumber  string
    Email        string
    Gender       string    `gorm:"type:varchar(10)"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}