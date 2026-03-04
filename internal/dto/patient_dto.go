package dto

// Request DTOs

type SearchPatientRequest struct {
    NationalID   string `form:"national_id"`
    PassportID   string `form:"passport_id"`
    FirstName    string `form:"first_name"`
    MiddleName   string `form:"middle_name"`
    LastName     string `form:"last_name"`
    DateOfBirth  string `form:"date_of_birth"`
    PhoneNumber  string `form:"phone_number"`
    Email        string `form:"email"`
    Page         int    `form:"page"`
    PageSize     int    `form:"page_size"`
}

// Response DTOs

type PatientResponse struct {
    ID           string  `json:"id"`
    HospitalID   string  `json:"hospital_id"`
    FirstNameTH  string  `json:"first_name_th"`
    MiddleNameTH string  `json:"middle_name_th"`
    LastNameTH   string  `json:"last_name_th"`
    FirstNameEN  string  `json:"first_name_en"`
    MiddleNameEN string  `json:"middle_name_en"`
    LastNameEN   string  `json:"last_name_en"`
    DateOfBirth  string  `json:"date_of_birth"`
    PatientHN    string  `json:"patient_hn"`
    NationalID   *string `json:"national_id"`
    PassportID   *string `json:"passport_id"`
    PhoneNumber  string  `json:"phone_number"`
    Email        string  `json:"email"`
    Gender       string  `json:"gender"`
}