package repository

import (
    "agnos-backend/internal/dto"
    "agnos-backend/internal/model"
    "gorm.io/gorm"
)

// Interface

type PatientRepository interface {
    Search(hospitalID string, req dto.SearchPatientRequest) ([]model.Patient, int64, error)
}

// Implementation

type patientRepo struct {
    db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
    return &patientRepo{db: db}
}

func (r *patientRepo) Search(hospitalID string, req dto.SearchPatientRequest) ([]model.Patient, int64, error) {
    var patients []model.Patient
    var total int64

    query := r.db.Model(&model.Patient{}).
        Where("hospital_id = ?", hospitalID)

    // Apply optional filters
    if req.NationalID != "" {
        query = query.Where("national_id = ?", req.NationalID)
    }
    if req.PassportID != "" {
        query = query.Where("passport_id = ?", req.PassportID)
    }
    if req.FirstName != "" {
        query = query.Where("first_name_th LIKE ? OR first_name_en LIKE ?",
            "%"+req.FirstName+"%", "%"+req.FirstName+"%")
    }
    if req.MiddleName != "" {
        query = query.Where("middle_name_th LIKE ? OR middle_name_en LIKE ?",
            "%"+req.MiddleName+"%", "%"+req.MiddleName+"%")
    }
    if req.LastName != "" {
        query = query.Where("last_name_th LIKE ? OR last_name_en LIKE ?",
            "%"+req.LastName+"%", "%"+req.LastName+"%")
    }
    if req.DateOfBirth != "" {
        query = query.Where("date_of_birth = ?", req.DateOfBirth)
    }
    if req.PhoneNumber != "" {
        query = query.Where("phone_number = ?", req.PhoneNumber)
    }
    if req.Email != "" {
        query = query.Where("email = ?", req.Email)
    }

    // Count total before pagination
    query.Count(&total)

    // Apply pagination
    page := req.Page
    pageSize := req.PageSize
    if page <= 0 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = 10
    }
    offset := (page - 1) * pageSize

    err := query.
        Limit(pageSize).
        Offset(offset).
        Find(&patients).Error

    return patients, total, err
}