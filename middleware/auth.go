package middleware

import (
	"net/http"
	"strconv"
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

		var staffID uint = claims.StaffID
		var hospitalID uint = claims.HospitalID
		if staffID == 0 {
			var mc jwt.MapClaims
			if _, err2 := jwt.ParseWithClaims(tokenString, &mc, func(token *jwt.Token) (interface{}, error) {
				return []byte(conf.JwtSecret), nil
			}); err2 == nil {
				readUint := func(key string) (uint, bool) {
					if v, ok := mc[key]; ok && v != nil {
						switch t := v.(type) {
						case float64:
							return uint(t), true
						case int64:
							return uint(t), true
						case int:
							return uint(t), true
						case string:
							if u, err := strconv.ParseUint(t, 10, 64); err == nil {
								return uint(u), true
							}
						}
					}
					return 0, false
				}

				if v, ok := readUint("staff_id"); ok {
					staffID = v
				} else if v, ok := readUint("StaffID"); ok {
					staffID = v
				} else if v, ok := readUint("staffId"); ok {
					staffID = v
				}
				if v, ok := readUint("hospital_id"); ok {
					hospitalID = v
				} else if v, ok := readUint("HospitalID"); ok {
					hospitalID = v
				} else if v, ok := readUint("hospitalId"); ok {
					hospitalID = v
				}
			}
		}

		if _, err := staffRepo.GetByID(staffID); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Staff not found"})
			return
		}

		if staffID != claims.StaffID || hospitalID != claims.HospitalID {
			c.Set(string(StaffContextKey), &StaffClaims{StaffID: staffID, HospitalID: hospitalID})
		} else {
			c.Set(string(StaffContextKey), claims)
		}
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
