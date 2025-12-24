package tests

import (
	"errors"
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
		a := "A"
		return []models.Patient{{ID: 1, FirstNameTH: &a}}, nil
	}}

	ph := handlers.NewPatientHandler(mock)
	r := gin.New()
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

func TestPatientGetByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockPatientService{GetByFn: func(hospitalID uint, id string) (*models.Patient, error) {
		a := "A"
		return &models.Patient{ID: 42, FirstNameTH: &a}, nil
	}}

	ph := handlers.NewPatientHandler(mock)
	r := gin.New()
	r.GET("/api/patient/:id", func(c *gin.Context) {
		c.Set("hospital_id", uint(2))
		ph.GetByID(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/patient/X", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestPatientGetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockPatientService{GetByFn: func(hospitalID uint, id string) (*models.Patient, error) {
		return nil, errors.New("not found")
	}}

	ph := handlers.NewPatientHandler(mock)
	r := gin.New()
	r.GET("/api/patient/:id", func(c *gin.Context) {
		c.Set("hospital_id", uint(2))
		ph.GetByID(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/patient/X", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestPatientGetByID_NoHospitalContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockPatientService{}
	ph := handlers.NewPatientHandler(mock)
	r := gin.New()
	r.GET("/api/patient/:id", func(c *gin.Context) { ph.GetByID(c) })

	req := httptest.NewRequest(http.MethodGet, "/api/patient/X", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPatientSearch_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockPatientService{SearchFn: func(hospitalID uint, filters map[string]interface{}) ([]models.Patient, error) {
		return nil, errors.New("boom")
	}}

	ph := handlers.NewPatientHandler(mock)
	r := gin.New()
	r.GET("/api/patient/search", func(c *gin.Context) {
		c.Set(string(middleware.StaffContextKey), &middleware.StaffClaims{HospitalID: 2})
		ph.Search(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/patient/search?national_id=X", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
