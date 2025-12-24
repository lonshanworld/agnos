package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"agnos_candidate_assignment/handlers"
	"agnos_candidate_assignment/middleware"
	"agnos_candidate_assignment/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type mockPatientService struct {
	SearchFn func(hospitalID uint, filters map[string]interface{}) ([]models.Patient, error)
	GetByFn  func(hospitalID uint, id string) (*models.Patient, error)
}

func (m *mockPatientService) Search(hospitalID uint, filters map[string]interface{}) ([]models.Patient, error) {
	return m.SearchFn(hospitalID, filters)
}
func (m *mockPatientService) GetByNationalOrPassport(hospitalID uint, id string) (*models.Patient, error) {
	return m.GetByFn(hospitalID, id)
}

func TestPatientSearch_Authorized_Positive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockPatientService{SearchFn: func(hospitalID uint, filters map[string]interface{}) ([]models.Patient, error) {
		return []models.Patient{{ID: 1, FirstNameTH: "A"}}, nil
	}}

	ph := handlers.NewPatientHandler(mock)
	r := gin.New()
	// attach middleware that injects claims
	r.GET("/api/patient/search", func(c *gin.Context) {
		c.Set(string(middleware.StaffContextKey), &middleware.StaffClaims{HospitalID: 2})
		ph.Search(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/patient/search?national_id=X", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestPatientSearch_Unauthorized_NoClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockPatientService{}
	ph := handlers.NewPatientHandler(mock)
	r := gin.New()
	r.GET("/api/patient/search", func(c *gin.Context) { ph.Search(c) })

	req := httptest.NewRequest(http.MethodGet, "/api/patient/search?national_id=X", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}
