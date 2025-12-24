package handlers

import (
	"net/http"

	"agnos_candidate_assignment/middleware"
	"agnos_candidate_assignment/services"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	patientService services.PatientServiceInterface
}

func NewPatientHandler(patientService services.PatientServiceInterface) *PatientHandler {
	return &PatientHandler{patientService: patientService}
}

func (patientHandler *PatientHandler) Search(c *gin.Context) {
	claims := middleware.GetStaffClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing staff claims"})
		return
	}
	hospitalID := claims.HospitalID
	filters := map[string]any{}
	if v := c.Query("national_id"); v != "" {
		filters["national_id"] = v
	}
	if v := c.Query("passport_id"); v != "" {
		filters["passport_id"] = v
	}
	if v := c.Query("first_name"); v != "" {
		filters["first_name"] = v
	}
	if v := c.Query("middle_name"); v != "" {
		filters["middle_name"] = v
	}
	if v := c.Query("last_name"); v != "" {
		filters["last_name"] = v
	}
	if v := c.Query("first_name_th"); v != "" {
		filters["first_name_th"] = v
	}
	if v := c.Query("middle_name_th"); v != "" {
		filters["middle_name_th"] = v
	}
	if v := c.Query("last_name_th"); v != "" {
		filters["last_name_th"] = v
	}
	if v := c.Query("first_name_en"); v != "" {
		filters["first_name_en"] = v
	}
	if v := c.Query("middle_name_en"); v != "" {
		filters["middle_name_en"] = v
	}
	if v := c.Query("last_name_en"); v != "" {
		filters["last_name_en"] = v
	}
	if v := c.Query("date_of_birth"); v != "" {
		filters["date_of_birth"] = v
	}
	if v := c.Query("phone_number"); v != "" {
		filters["phone_number"] = v
	}
	if v := c.Query("email"); v != "" {
		filters["email"] = v
	}

	results, err := patientHandler.patientService.Search(hospitalID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"patients": results})
}

func (h *PatientHandler) GetByID(c *gin.Context) {
	raw, ok := c.Get("hospital_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hospital context missing"})
		return
	}
	hospitalID := raw.(uint)
	id := c.Param("id")
	p, err := h.patientService.GetByNationalOrPassport(hospitalID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}
