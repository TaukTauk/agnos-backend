package middleware

import (
    "net/http"
    "os"
    "sync"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// JWT Blacklist

var (
    blacklist   = make(map[string]bool)
    blacklistMu sync.RWMutex
)

func BlacklistToken(token string) {
    blacklistMu.Lock()
    defer blacklistMu.Unlock()
    blacklist[token] = true
}

func isBlacklisted(token string) bool {
    blacklistMu.RLock()
    defer blacklistMu.RUnlock()
    return blacklist[token]
}

// API Key Middleware

func APIKeyAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")

        if apiKey == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "missing API key",
            })
            c.Abort()
            return
        }

        if apiKey != os.Getenv("API_KEY") {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "invalid API key",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// JWT Middleware

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString, err := c.Cookie("jwt")
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "missing token",
            })
            c.Abort()
            return
        }

        // Check blacklist
        if isBlacklisted(tokenString) {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "token has been invalidated",
            })
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "invalid or expired token",
            })
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "invalid token claims",
            })
            c.Abort()
            return
        }

        // Store raw token string for blacklisting on logout
        c.Set("token",       tokenString)
        c.Set("staff_id",    claims["staff_id"])
        c.Set("hospital_id", claims["hospital_id"])
        c.Set("username",    claims["username"])

        c.Next()
    }
}