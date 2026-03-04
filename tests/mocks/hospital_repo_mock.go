package mocks

import (
    "agnos-backend/internal/model"
    "github.com/stretchr/testify/mock"
)

type MockHospitalRepository struct {
    mock.Mock
}

func (m *MockHospitalRepository) FindByCode(code string) (*model.Hospital, error) {
    args := m.Called(code)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.Hospital), args.Error(1)
}