package dto

import "time"

// Request DTOs

type CreateStaffRequest struct {
    Username     string `json:"username"      binding:"required,min=3,max=50"`
    Password     string `json:"password"      binding:"required,min=8,max=72"`
    HospitalCode string `json:"hospital_code" binding:"required"`
}

type LoginStaffRequest struct {
    Username      string `json:"username"       binding:"required"`
    Password      string `json:"password"       binding:"required"`
    HospitalCode  string `json:"hospital_code"  binding:"required"`
}

// Response DTOs

type StaffResponse struct {
    ID          string    `json:"id"`
    Username    string    `json:"username"`
    HospitalID  string    `json:"hospital_id"`
    CreatedAt   time.Time `json:"created_at"`
}

type LoginResponse struct {
    StaffID     string `json:"staff_id"`
    Username    string `json:"username"`
    HospitalID  string `json:"hospital_id"`
}