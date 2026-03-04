package repository

import (
    "agnos-backend/internal/model"
    "gorm.io/gorm"
)

// Interface

type HospitalRepository interface {
    FindByCode(code string) (*model.Hospital, error)
}

// Implementation

type hospitalRepo struct {
    db *gorm.DB
}

func NewHospitalRepository(db *gorm.DB) HospitalRepository {
    return &hospitalRepo{db: db}
}

func (r *hospitalRepo) FindByCode(code string) (*model.Hospital, error) {
    var hospital model.Hospital
    err := r.db.Where("code = ?", code).First(&hospital).Error
    if err != nil {
        return nil, err
    }
    return &hospital, nil
}