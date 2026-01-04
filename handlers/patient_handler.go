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

// Search godoc
// @Summary      Search patients
// @Description  Search for patients by various criteria (requires authentication)
// @Tags         patients
// @Accept       json
// @Produce      json
// @Param        national_id query string false "National ID"
// @Param        passport_id query string false "Passport ID"
// @Param        first_name query string false "First name (any language)"
// @Param        middle_name query string false "Middle name (any language)"
// @Param        last_name query string false "Last name (any language)"
// @Param        first_name_th query string false "First name (Thai)"
// @Param        middle_name_th query string false "Middle name (Thai)"
// @Param        last_name_th query string false "Last name (Thai)"
// @Param        first_name_en query string false "First name (English)"
// @Param        middle_name_en query string false "Middle name (English)"
// @Param        last_name_en query string false "Last name (English)"
// @Param        date_of_birth query string false "Date of birth"
// @Param        phone_number query string false "Phone number"
// @Param        email query string false "Email"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /patient/search [get]
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

// GetByID godoc
// @Summary      Get patient by ID
// @Description  Retrieve patient by national ID or passport ID (public endpoint)
// @Tags         patients
// @Accept       json
// @Produce      json
// @Param        hospital path string true "Hospital name"
// @Param        id path string true "National ID or Passport ID"
// @Success      200  {object}  models.Patient
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /:hospital/patient/search/{id} [get]
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
