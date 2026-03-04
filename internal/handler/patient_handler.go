package handler

import (
	"log"
    "net/http"

    "github.com/gin-gonic/gin"

    "agnos-backend/internal/dto"
    "agnos-backend/internal/service"
)

// Interface

type PatientHandler interface {
    Search(c *gin.Context)
}

// Implementation

type patientHandler struct {
    patientService service.PatientService
}

func NewPatientHandler(patientService service.PatientService) PatientHandler {
    return &patientHandler{patientService: patientService}
}

// Search handles GET /patient/search
func (h *patientHandler) Search(c *gin.Context) {
    var req dto.SearchPatientRequest

    // Parse query params
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Success: false,
            Error:   err.Error(),
        })
        return
    }

    // Get hospital_id injected by JWT middleware
    hospitalID := c.GetString("hospital_id")
	staffID    := c.GetString("staff_id")

    // Call service
    result, err := h.patientService.Search(hospitalID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Success: false,
            Error:   "internal server error",
        })
        return
    }

	// Audit log
    log.Printf("[AUDIT] staff_id=%s hospital_id=%s action=patient_search filters=%+v",
        staffID,
        hospitalID,
        req,
    )

    c.JSON(http.StatusOK, result)
}