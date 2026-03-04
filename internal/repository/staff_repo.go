package repository

import (
    "agnos-backend/internal/model"
    "gorm.io/gorm"
)

// Interface

type StaffRepository interface {
    Create(staff *model.Staff) error
    FindByUsernameAndHospital(username string, hospitalID string) (*model.Staff, error)
}

// Implementation

type staffRepo struct {
    db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) StaffRepository {
    return &staffRepo{db: db}
}

func (r *staffRepo) Create(staff *model.Staff) error {
    return r.db.Create(staff).Error
}

func (r *staffRepo) FindByUsernameAndHospital(username string, hospitalID string) (*model.Staff, error) {
    var staff model.Staff
    err := r.db.
        Where("username = ? AND hospital_id = ?", username, hospitalID).
        First(&staff).Error
    if err != nil {
        return nil, err
    }
    return &staff, nil
}