package middleware

import (
	"net/http"
	"strings"

	"agnos_candidate_assignment/config"
	"agnos_candidate_assignment/repositories"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type StaffClaims struct {
	StaffID    uint `json:"staff_id"`
	HospitalID uint `json:"hospital_id"`
	jwt.RegisteredClaims
}

type ContextStaffKey string

const StaffContextKey ContextStaffKey = "staff_claims"

func JWTAuth(conf *config.Config, staffRepo *repositories.StaffRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.Fields(auth)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := parts[1]
		claims := &StaffClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.JwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if _, err := staffRepo.GetByID(claims.StaffID); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Staff not found"})
			return
		}

		c.Set(string(StaffContextKey), claims)
		c.Next()
	}
}

func GetStaffClaims(c *gin.Context) *StaffClaims {
	if v, ok := c.Get(string(StaffContextKey)); ok {
		if claims, ok := v.(*StaffClaims); ok {
			return claims
		}
	}
	return nil
}
