package service

import (
    "errors"
    "os"
    "time"
	"strconv"
	"unicode"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    "agnos-backend/internal/dto"
    "agnos-backend/internal/model"
    "agnos-backend/internal/repository"
)

// Interface

type StaffService interface {
    CreateStaff(req dto.CreateStaffRequest) (*dto.StaffResponse, error)
    Login(req dto.LoginStaffRequest) (*dto.LoginResponse, string, error)
}

// Implementation

type staffService struct {
    staffRepo    repository.StaffRepository
    hospitalRepo repository.HospitalRepository
}

func NewStaffService(
    staffRepo repository.StaffRepository,
    hospitalRepo repository.HospitalRepository,
) StaffService {
    return &staffService{
        staffRepo:    staffRepo,
        hospitalRepo: hospitalRepo,
    }
}

// CreateStaff creates a new staff member
func (s *staffService) CreateStaff(req dto.CreateStaffRequest) (*dto.StaffResponse, error) {
    // 1. Find hospital by code
    hospital, err := s.hospitalRepo.FindByCode(req.HospitalCode)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("hospital not found")
        }
        return nil, err
    }

	// Validate username
	if err := validateUsername(req.Username); err != nil {
	    return nil, err
	}

	// Validate password
	if err := validatePassword(req.Password); err != nil {
	    return nil, err
	}

    // 2. Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errors.New("failed to hash password")
    }

    // 3. Build the staff model
    staff := &model.Staff{
        ID:           uuid.New().String(),
        Username:     req.Username,
        PasswordHash: string(hashedPassword),
        HospitalID:   hospital.ID,
    }

    // 4. Save to DB
    err = s.staffRepo.Create(staff)
    if err != nil {
        return nil, errors.New("username already exists")
    }

    // 5. Return response DTO
    return &dto.StaffResponse{
        ID:         staff.ID,
        Username:   staff.Username,
        HospitalID: staff.HospitalID,
        CreatedAt:  staff.CreatedAt,
    }, nil
}

// Login validates credentials and returns a JWT token
func (s *staffService) Login(req dto.LoginStaffRequest) (*dto.LoginResponse, string, error) {
    // 1. Find hospital by code
    hospital, err := s.hospitalRepo.FindByCode(req.HospitalCode)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, "", errors.New("hospital not found")
        }
        return nil, "", err
    }

    // 2. Find staff by username + hospital
    staff, err := s.staffRepo.FindByUsernameAndHospital(req.Username, hospital.ID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, "", errors.New("invalid credentials")
        }
        return nil, "", err
    }

    // 3. Compare password
    err = bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password))
    if err != nil {
        return nil, "", errors.New("invalid credentials")
    }

    // 4. Generate JWT
    token, err := generateJWT(staff.ID, staff.HospitalID, staff.Username)
    if err != nil {
        return nil, "", errors.New("failed to generate token")
    }

    // 5. Return response DTO + token
    return &dto.LoginResponse{
        StaffID:    staff.ID,
        Username:   staff.Username,
        HospitalID: staff.HospitalID,
    }, token, nil
}

// generateJWT creates a signed JWT token
func generateJWT(staffID string, hospitalID string, username string) (string, error) {
	expiryHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if err != nil || expiryHours <= 0 {
		expiryHours = 24 // fallback default
	}

    claims := jwt.MapClaims{
        "staff_id":    staffID,
        "hospital_id": hospitalID,
        "username":    username,
        "exp":         time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// helper function
func validatePassword(password string) error {
    var hasUpper, hasLower, hasNumber, hasSpecial bool
    for _, c := range password {
        switch {
        case unicode.IsUpper(c):
            hasUpper = true
        case unicode.IsLower(c):
            hasLower = true
        case unicode.IsNumber(c):
            hasNumber = true
        case unicode.IsPunct(c) || unicode.IsSymbol(c):
            hasSpecial = true
        }
    }
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    if !hasUpper {
        return errors.New("password must contain at least one uppercase letter")
    }
    if !hasLower {
        return errors.New("password must contain at least one lowercase letter")
    }
    if !hasNumber {
        return errors.New("password must contain at least one number")
    }
    if !hasSpecial {
        return errors.New("password must contain at least one special character")
    }
    return nil
}

func validateUsername(username string) error {
    if len(username) < 3 || len(username) > 50 {
        return errors.New("username must be between 3 and 50 characters")
    }
    for _, c := range username {
        if !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '.' && c != '_' {
            return errors.New("username can only contain letters, numbers, dots and underscores")
        }
    }
    return nil
}