package service_test

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"

    "agnos-backend/internal/dto"
    "agnos-backend/internal/model"
    "agnos-backend/internal/service"
    "agnos-backend/tests/mocks"
)

// Search Tests

func TestSearchPatient_Success(t *testing.T) {
    mockPatientRepo := new(mocks.MockPatientRepository)

    nationalID := "1100100011001"
    mockPatientRepo.On("Search", "hospital-a-uuid", dto.SearchPatientRequest{
        FirstName: "Somchai",
        Page:      1,
        PageSize:  10,
    }).Return([]model.Patient{
        {
            ID:          "patient-uuid-1",
            HospitalID:  "hospital-a-uuid",
            FirstNameEN: "Somchai",
            LastNameEN:  "Jaidee",
            NationalID:  &nationalID,
            Gender:      "male",
        },
    }, int64(1), nil)

    svc := service.NewPatientService(mockPatientRepo)

    result, err := svc.Search("hospital-a-uuid", dto.SearchPatientRequest{
        FirstName: "Somchai",
        Page:      1,
        PageSize:  10,
    })

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, int64(1), result.Pagination.Total)
    assert.Equal(t, 1, result.Pagination.TotalPages)

    patients := result.Data.([]dto.PatientResponse)
    assert.Len(t, patients, 1)
    assert.Equal(t, "Somchai", patients[0].FirstNameEN)
    mockPatientRepo.AssertExpectations(t)
}

func TestSearchPatient_NoResults(t *testing.T) {
    mockPatientRepo := new(mocks.MockPatientRepository)

    mockPatientRepo.On("Search", "hospital-a-uuid", dto.SearchPatientRequest{
        FirstName: "Unknown",
        Page:      1,
        PageSize:  10,
    }).Return([]model.Patient{}, int64(0), nil)

    svc := service.NewPatientService(mockPatientRepo)

    result, err := svc.Search("hospital-a-uuid", dto.SearchPatientRequest{
        FirstName: "Unknown",
        Page:      1,
        PageSize:  10,
    })

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, int64(0), result.Pagination.Total)
    assert.Nil(t, result.Data)
    mockPatientRepo.AssertExpectations(t)
}

func TestSearchPatient_DefaultPagination(t *testing.T) {
    mockPatientRepo := new(mocks.MockPatientRepository)

    // Page and PageSize should default to 1 and 10
    mockPatientRepo.On("Search", "hospital-a-uuid", dto.SearchPatientRequest{
        Page:     1,
        PageSize: 10,
    }).Return([]model.Patient{}, int64(0), nil)

    svc := service.NewPatientService(mockPatientRepo)

    // Send with no page/page_size
    result, err := svc.Search("hospital-a-uuid", dto.SearchPatientRequest{})

    assert.NoError(t, err)
    assert.Equal(t, 1, result.Pagination.Page)
    assert.Equal(t, 10, result.Pagination.PageSize)
    mockPatientRepo.AssertExpectations(t)
}

func TestSearchPatient_RepositoryError(t *testing.T) {
    mockPatientRepo := new(mocks.MockPatientRepository)

    mockPatientRepo.On("Search", "hospital-a-uuid", dto.SearchPatientRequest{
        Page:     1,
        PageSize: 10,
    }).Return([]model.Patient{}, int64(0), errors.New("db connection error"))

    svc := service.NewPatientService(mockPatientRepo)

    result, err := svc.Search("hospital-a-uuid", dto.SearchPatientRequest{})

    assert.Nil(t, result)
    assert.EqualError(t, err, "db connection error")
    mockPatientRepo.AssertExpectations(t)
}

func TestSearchPatient_PaginationCalculation(t *testing.T) {
    mockPatientRepo := new(mocks.MockPatientRepository)

    mockPatientRepo.On("Search", "hospital-a-uuid", dto.SearchPatientRequest{
        Page:     1,
        PageSize: 3,
    }).Return([]model.Patient{
        {ID: "p1"}, {ID: "p2"}, {ID: "p3"},
    }, int64(7), nil) // 7 total → 3 pages (ceil(7/3))

    svc := service.NewPatientService(mockPatientRepo)

    result, err := svc.Search("hospital-a-uuid", dto.SearchPatientRequest{
        Page:     1,
        PageSize: 3,
    })

    assert.NoError(t, err)
    assert.Equal(t, int64(7), result.Pagination.Total)
    assert.Equal(t, 3, result.Pagination.TotalPages) // ceil(7/3) = 3
    mockPatientRepo.AssertExpectations(t)
}