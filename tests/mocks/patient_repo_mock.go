package mocks

import (
    "agnos-backend/internal/dto"
    "agnos-backend/internal/model"
    "github.com/stretchr/testify/mock"
)

type MockPatientRepository struct {
    mock.Mock
}

func (m *MockPatientRepository) Search(hospitalID string, req dto.SearchPatientRequest) ([]model.Patient, int64, error) {
    args := m.Called(hospitalID, req)
    return args.Get(0).([]model.Patient), args.Get(1).(int64), args.Error(2)
}