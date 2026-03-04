package service_test

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "gorm.io/gorm"

    "agnos-backend/internal/dto"
    "agnos-backend/internal/model"
    "agnos-backend/internal/service"
    "agnos-backend/tests/mocks"
)

// CreateStaff Tests

func TestCreateStaff_Success(t *testing.T) {
    // Arrange
    mockHospitalRepo := new(mocks.MockHospitalRepository)
    mockStaffRepo    := new(mocks.MockStaffRepository)

    mockHospitalRepo.On("FindByCode", "HOSPITAL_A").
        Return(&model.Hospital{
            ID:   "hospital-a-uuid",
            Name: "Hospital A",
            Code: "HOSPITAL_A",
        }, nil)

    mockStaffRepo.On("Create", mock.AnythingOfType("*model.Staff")).
        Return(nil)

    svc := service.NewStaffService(mockStaffRepo, mockHospitalRepo)

    req := dto.CreateStaffRequest{
        Username:     "newstaff",
        Password:     "Password123!",
        HospitalCode: "HOSPITAL_A",
    }

    // Act
    result, err := svc.CreateStaff(req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "newstaff", result.Username)
    assert.Equal(t, "hospital-a-uuid", result.HospitalID)
    mockHospitalRepo.AssertExpectations(t)
    mockStaffRepo.AssertExpectations(t)
}

func TestCreateStaff_HospitalNotFound(t *testing.T) {
    mockHospitalRepo := new(mocks.MockHospitalRepository)
    mockStaffRepo    := new(mocks.MockStaffRepository)

    mockHospitalRepo.On("FindByCode", "INVALID_HOSPITAL").
        Return(nil, gorm.ErrRecordNotFound)

    svc := service.NewStaffService(mockStaffRepo, mockHospitalRepo)

    req := dto.CreateStaffRequest{
        Username:     "newstaff",
        Password:     "Password123!",
        HospitalCode: "INVALID_HOSPITAL",
    }

    result, err := svc.CreateStaff(req)

    assert.Nil(t, result)
    assert.EqualError(t, err, "hospital not found")
    mockHospitalRepo.AssertExpectations(t)
}

func TestCreateStaff_DuplicateUsername(t *testing.T) {
    mockHospitalRepo := new(mocks.MockHospitalRepository)
    mockStaffRepo    := new(mocks.MockStaffRepository)

    mockHospitalRepo.On("FindByCode", "HOSPITAL_A").
        Return(&model.Hospital{
            ID:   "hospital-a-uuid",
            Code: "HOSPITAL_A",
        }, nil)

    mockStaffRepo.On("Create", mock.AnythingOfType("*model.Staff")).
        Return(errors.New("duplicate key"))

    svc := service.NewStaffService(mockStaffRepo, mockHospitalRepo)

    req := dto.CreateStaffRequest{
        Username:     "existingstaff",
        Password:     "Password123!",
        HospitalCode: "HOSPITAL_A",
    }

    result, err := svc.CreateStaff(req)

    assert.Nil(t, result)
    assert.EqualError(t, err, "username already exists")
    mockHospitalRepo.AssertExpectations(t)
    mockStaffRepo.AssertExpectations(t)
}

// Login Tests

func TestLogin_Success(t *testing.T) {
    // Need JWT_SECRET set for token generation
    t.Setenv("JWT_SECRET", "test-secret")
    t.Setenv("JWT_EXPIRY_HOURS", "24")

    mockHospitalRepo := new(mocks.MockHospitalRepository)
    mockStaffRepo    := new(mocks.MockStaffRepository)

    mockHospitalRepo.On("FindByCode", "HOSPITAL_A").
        Return(&model.Hospital{
            ID:   "hospital-a-uuid",
            Code: "HOSPITAL_A",
        }, nil)

    // Use bcrypt hash of "Password123!"
    mockStaffRepo.On("FindByUsernameAndHospital", "staff.hospital.a", "hospital-a-uuid").
        Return(&model.Staff{
            ID:           "staff-uuid",
            Username:     "staff.hospital.a",
            PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
            HospitalID:   "hospital-a-uuid",
        }, nil)

    svc := service.NewStaffService(mockStaffRepo, mockHospitalRepo)

    req := dto.LoginStaffRequest{
        Username:     "staff.hospital.a",
        Password:     "password",
        HospitalCode: "HOSPITAL_A",
    }

    result, token, err := svc.Login(req)

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.NotEmpty(t, token)
    assert.Equal(t, "staff.hospital.a", result.Username)
    mockHospitalRepo.AssertExpectations(t)
    mockStaffRepo.AssertExpectations(t)
}

func TestLogin_HospitalNotFound(t *testing.T) {
    mockHospitalRepo := new(mocks.MockHospitalRepository)
    mockStaffRepo    := new(mocks.MockStaffRepository)

    mockHospitalRepo.On("FindByCode", "INVALID").
        Return(nil, gorm.ErrRecordNotFound)

    svc := service.NewStaffService(mockStaffRepo, mockHospitalRepo)

    req := dto.LoginStaffRequest{
        Username:     "staff.hospital.a",
        Password:     "password",
        HospitalCode: "INVALID",
    }

    result, token, err := svc.Login(req)

    assert.Nil(t, result)
    assert.Empty(t, token)
    assert.EqualError(t, err, "hospital not found")
}

func TestLogin_StaffNotFound(t *testing.T) {
    mockHospitalRepo := new(mocks.MockHospitalRepository)
    mockStaffRepo    := new(mocks.MockStaffRepository)

    mockHospitalRepo.On("FindByCode", "HOSPITAL_A").
        Return(&model.Hospital{
            ID:   "hospital-a-uuid",
            Code: "HOSPITAL_A",
        }, nil)

    mockStaffRepo.On("FindByUsernameAndHospital", "nobody", "hospital-a-uuid").
        Return(nil, gorm.ErrRecordNotFound)

    svc := service.NewStaffService(mockStaffRepo, mockHospitalRepo)

    req := dto.LoginStaffRequest{
        Username:     "nobody",
        Password:     "password",
        HospitalCode: "HOSPITAL_A",
    }

    result, token, err := svc.Login(req)

    assert.Nil(t, result)
    assert.Empty(t, token)
    assert.EqualError(t, err, "invalid credentials")
}

func TestLogin_WrongPassword(t *testing.T) {
    mockHospitalRepo := new(mocks.MockHospitalRepository)
    mockStaffRepo    := new(mocks.MockStaffRepository)

    mockHospitalRepo.On("FindByCode", "HOSPITAL_A").
        Return(&model.Hospital{
            ID:   "hospital-a-uuid",
            Code: "HOSPITAL_A",
        }, nil)

    mockStaffRepo.On("FindByUsernameAndHospital", "staff.hospital.a", "hospital-a-uuid").
        Return(&model.Staff{
            ID:           "staff-uuid",
            Username:     "staff.hospital.a",
            PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
            HospitalID:   "hospital-a-uuid",
        }, nil)

    svc := service.NewStaffService(mockStaffRepo, mockHospitalRepo)

    req := dto.LoginStaffRequest{
        Username:     "staff.hospital.a",
        Password:     "wrongpassword",
        HospitalCode: "HOSPITAL_A",
    }

    result, token, err := svc.Login(req)

    assert.Nil(t, result)
    assert.Empty(t, token)
    assert.EqualError(t, err, "invalid credentials")
}

func TestCreateStaff_WeakPassword(t *testing.T) {
    mockHospitalRepo := new(mocks.MockHospitalRepository)
    mockStaffRepo    := new(mocks.MockStaffRepository)

    mockHospitalRepo.On("FindByCode", "HOSPITAL_A").
        Return(&model.Hospital{
            ID:   "hospital-a-uuid",
            Code: "HOSPITAL_A",
        }, nil)

    svc := service.NewStaffService(mockStaffRepo, mockHospitalRepo)

    req := dto.CreateStaffRequest{
        Username:     "newstaff",
        Password:     "weak",
        HospitalCode: "HOSPITAL_A",
    }

    result, err := svc.CreateStaff(req)

    assert.Nil(t, result)
    assert.Error(t, err)
}