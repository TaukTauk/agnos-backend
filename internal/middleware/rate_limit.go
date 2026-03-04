package middleware

import (
    "github.com/gin-gonic/gin"
    ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
    "github.com/ulule/limiter/v3"
    "github.com/ulule/limiter/v3/drivers/store/memory"
    "time"
)

func RateLimiter(limit int64) gin.HandlerFunc {
    // limit requests per minute per IP
    rate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  limit,
    }

    store := memory.NewStore()
    instance := limiter.New(store, rate)

    return ginLimiter.NewMiddleware(instance)
}