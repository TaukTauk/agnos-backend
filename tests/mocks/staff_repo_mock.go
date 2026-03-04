package mocks

import (
    "agnos-backend/internal/model"
    "github.com/stretchr/testify/mock"
)

type MockStaffRepository struct {
    mock.Mock
}

func (m *MockStaffRepository) Create(staff *model.Staff) error {
    args := m.Called(staff)
    return args.Error(0)
}

func (m *MockStaffRepository) FindByUsernameAndHospital(username string, hospitalID string) (*model.Staff, error) {
    args := m.Called(username, hospitalID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.Staff), args.Error(1)
}