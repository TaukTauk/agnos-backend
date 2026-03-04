package handler

import (
    "net/http"
	"strconv"
	"os"

    "github.com/gin-gonic/gin"

    "agnos-backend/internal/dto"
    "agnos-backend/internal/service"
	"agnos-backend/internal/middleware"
)

// Interface

type StaffHandler interface {
    Create(c *gin.Context)
    Login(c *gin.Context)
    Logout(c *gin.Context)
}

// Implementation

type staffHandler struct {
    staffService service.StaffService
}

func NewStaffHandler(staffService service.StaffService) StaffHandler {
    return &staffHandler{staffService: staffService}
}

// Create handles POST /staff/create
func (h *staffHandler) Create(c *gin.Context) {
    var req dto.CreateStaffRequest

    // Parse and validate request body
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Success: false,
            Error:   err.Error(),
        })
        return
    }

    // Call service
    result, err := h.staffService.CreateStaff(req)
    if err != nil {
        switch err.Error() {
        case "hospital not found":
            c.JSON(http.StatusNotFound, dto.ErrorResponse{
                Success: false,
                Error:   err.Error(),
            })
        case "username already exists":
            c.JSON(http.StatusConflict, dto.ErrorResponse{
                Success: false,
                Error:   err.Error(),
            })
        default:
            c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
                Success: false,
                Error:   "internal server error",
            })
        }
        return
    }

    c.JSON(http.StatusCreated, dto.SuccessResponse{
        Success: true,
        Data:    result,
    })
}

// Login handles POST /staff/login
func (h *staffHandler) Login(c *gin.Context) {
    var req dto.LoginStaffRequest

    // Parse and validate request body
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Success: false,
            Error:   err.Error(),
        })
        return
    }

    // Call service
    result, token, err := h.staffService.Login(req)
    if err != nil {
        switch err.Error() {
        case "hospital not found":
            c.JSON(http.StatusNotFound, dto.ErrorResponse{
                Success: false,
                Error:   err.Error(),
            })
        case "invalid credentials":
            c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
                Success: false,
                Error:   err.Error(),
            })
        default:
            c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
                Success: false,
                Error:   "internal server error",
            })
        }
        return
    }

	expiryHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if err != nil || expiryHours <= 0 {
		expiryHours = 24 // fallback default
	}

    // Set JWT as httpOnly cookie
    c.SetCookie(
        "jwt",   // name
        token,   // value
        expiryHours * 3600,   // maxAge (hours in seconds)
        "/",     // path
        "",      // domain
        false,   // secure (set true in production with HTTPS)
        true,    // httpOnly
    )

    c.JSON(http.StatusOK, dto.SuccessResponse{
        Success: true,
        Data:    result,
    })
}

// Logout handles POST /staff/logout
func (h *staffHandler) Logout(c *gin.Context) {
    // Blacklist the current token
    token := c.GetString("token")
    if token != "" {
        middleware.BlacklistToken(token)
    }

    // Clear the cookie
    c.SetCookie("jwt", "", -1, "/", "", false, true)

    c.JSON(http.StatusOK, dto.SuccessResponse{
        Success: true,
        Data:    "logged out successfully",
    })
}