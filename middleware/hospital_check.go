package middleware

import (
	"net/http"
	"strconv"

	"agnos_candidate_assignment/repositories"

	"github.com/gin-gonic/gin"
)

func RequireHospitalMatch(hospRepo *repositories.HospitalRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		hParam := c.Param("hospital")
		var hospID uint
		if id, err := strconv.Atoi(hParam); err == nil {
			hospID = uint(id)
		} else {
			h, err := hospRepo.FindByName(hParam)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "hospital not found"})
				return
			}
			hospID = h.ID
		}

		claims := GetStaffClaims(c)
		if claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		if claims.HospitalID != hospID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "staff not allowed for this hospital"})
			return
		}
		c.Set("hospital_id", hospID)
		c.Next()
	}
}
