package main

import (
    "log"
    "net/http"
    "os"
    "strconv"
	"strings"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"

    "agnos-backend/config"
    "agnos-backend/internal/dto"
    "agnos-backend/internal/handler"
    "agnos-backend/internal/middleware"
    "agnos-backend/internal/repository"
    "agnos-backend/internal/service"
)

func getEnvInt(key string, defaultVal int64) int64 {
    val, err := strconv.ParseInt(os.Getenv(key), 10, 64)
    if err != nil || val <= 0 {
        return defaultVal
    }
    return val
}

func main() {
    // 1. Load .env
    config.Load()

    ginMode := os.Getenv("GIN_MODE")
    if ginMode == "" {
        ginMode = "debug"
    }
    gin.SetMode(ginMode)

    // 2. Connect to DB
    db, err := config.ConnectDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // 3. Run migrations
    config.Migrate(db)

    // 4. Wire up dependencies
    hospitalRepo := repository.NewHospitalRepository(db)
    staffRepo    := repository.NewStaffRepository(db)
    patientRepo  := repository.NewPatientRepository(db)

    staffService   := service.NewStaffService(staffRepo, hospitalRepo)
    patientService := service.NewPatientService(patientRepo)

    staffHandler   := handler.NewStaffHandler(staffService)
    patientHandler := handler.NewPatientHandler(patientService)

    // 5. Setup Gin router
    r := gin.Default()

    // CORS
	corsOrigins := os.Getenv("CORS_ORIGIN")
	allowedOrigins := []string{}
	if corsOrigins != "" {
	    for _, origin := range strings.Split(corsOrigins, ",") {
	        allowedOrigins = append(allowedOrigins, strings.TrimSpace(origin))
	    }
	}

	r.Use(cors.New(cors.Config{
	    AllowOrigins:     allowedOrigins,
	    AllowMethods:     []string{"GET", "POST"},
	    AllowHeaders:     []string{"Origin", "Content-Type", "X-API-Key"},
	    AllowCredentials: true,
	}))

    // 6. Global middleware — API key required for all routes
    r.Use(middleware.APIKeyAuth())

    // Rate limit values from .env with fallback defaults
    loginLimit       := getEnvInt("RATE_LIMIT_LOGIN", 5)
    createStaffLimit := getEnvInt("RATE_LIMIT_CREATE_STAFF", 10)
    searchLimit      := getEnvInt("RATE_LIMIT_PATIENT_SEARCH", 30)

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, dto.SuccessResponse{
            Success: true,
            Data:    "ok",
        })
    })

    // 7. Public routes (API key only)
    r.POST("/staff/create", middleware.RateLimiter(createStaffLimit), staffHandler.Create)
    r.POST("/staff/login",  middleware.RateLimiter(loginLimit),       staffHandler.Login)

    // 8. Protected routes (API key + JWT)
    protected := r.Group("/")
    protected.Use(middleware.JWTAuth())
    {
        protected.POST("/staff/logout",  staffHandler.Logout)
        protected.GET("/patient/search", middleware.RateLimiter(searchLimit), patientHandler.Search)
    }

    // 9. Start server
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server running on port %s", port)
    r.Run(":" + port)
}