package service

import (
    "math"
	"strings"

    "agnos-backend/internal/dto"
    "agnos-backend/internal/repository"
)

// Interface

type PatientService interface {
    Search(hospitalID string, req dto.SearchPatientRequest) (*dto.PaginatedResponse, error)
}

// Implementation

type patientService struct {
    patientRepo repository.PatientRepository
}

func NewPatientService(patientRepo repository.PatientRepository) PatientService {
    return &patientService{patientRepo: patientRepo}
}

func (s *patientService) Search(hospitalID string, req dto.SearchPatientRequest) (*dto.PaginatedResponse, error) {
	// Sanitize inputs
    req.FirstName   = strings.TrimSpace(req.FirstName)
    req.MiddleName  = strings.TrimSpace(req.MiddleName)
    req.LastName    = strings.TrimSpace(req.LastName)
    req.NationalID  = strings.TrimSpace(req.NationalID)
    req.PassportID  = strings.TrimSpace(req.PassportID)
    req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
    req.Email       = strings.TrimSpace(req.Email)

    // 1. Set defaults for pagination
    if req.Page <= 0 {
        req.Page = 1
    }
    if req.PageSize <= 0 {
        req.PageSize = 10
    }

    // 2. Call repository
    patients, total, err := s.patientRepo.Search(hospitalID, req)
    if err != nil {
        return nil, err
    }

    // 3. Map models to response DTOs
    var result []dto.PatientResponse
    for _, p := range patients {
        result = append(result, dto.PatientResponse{
            ID:           p.ID,
            HospitalID:   p.HospitalID,
            FirstNameTH:  p.FirstNameTH,
            MiddleNameTH: p.MiddleNameTH,
            LastNameTH:   p.LastNameTH,
            FirstNameEN:  p.FirstNameEN,
            MiddleNameEN: p.MiddleNameEN,
            LastNameEN:   p.LastNameEN,
            DateOfBirth:  p.DateOfBirth,
            PatientHN:    p.PatientHN,
            NationalID:   p.NationalID,
            PassportID:   p.PassportID,
            PhoneNumber:  p.PhoneNumber,
            Email:        p.Email,
            Gender:       p.Gender,
        })
    }

    // 4. Calculate total pages
    totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

    // 5. Return paginated response
    return &dto.PaginatedResponse{
        Success: true,
        Data:    result,
        Pagination: dto.Pagination{
            Page:       req.Page,
            PageSize:   req.PageSize,
            Total:      total,
            TotalPages: totalPages,
        },
    }, nil
}