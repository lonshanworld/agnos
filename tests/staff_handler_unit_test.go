package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"agnos_candidate_assignment/handlers"
	"agnos_candidate_assignment/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type mockAuthService struct {
	RegisterFn func(hospital, username, password string) (*models.Staff, error)
	LoginFn    func(hospital, username, password string) (string, *models.Staff, error)
}

func (m *mockAuthService) Register(hospital, username, password string) (*models.Staff, error) {
	return m.RegisterFn(hospital, username, password)
}
func (m *mockAuthService) Login(hospital, username, password string) (string, *models.Staff, error) {
	return m.LoginFn(hospital, username, password)
}

func TestStaffRegister_PositiveAndLogin_Positive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockAuthService{
		RegisterFn: func(hospital, username, password string) (*models.Staff, error) {
			return &models.Staff{ID: 7, UserName: username}, nil
		},
		LoginFn: func(hospital, username, password string) (string, *models.Staff, error) {
			return "tok", &models.Staff{ID: 7, UserName: username, HospitalID: 2}, nil
		},
	}

	sh := handlers.NewStaffHandler(mock)
	router := gin.New()
	router.POST("/api/:hospital/staff/create", sh.Register)
	router.POST("/api/:hospital/staff/login", sh.Login)

	reg := map[string]string{"username": "u1", "password": "p"}
	rb, _ := json.Marshal(reg)
	rreq := httptest.NewRequest(http.MethodPost, "/api/Hosp/staff/create", bytes.NewReader(rb))
	rreq.Header.Set("Content-Type", "application/json")
	rrec := httptest.NewRecorder()
	router.ServeHTTP(rrec, rreq)
	require.Equal(t, http.StatusCreated, rrec.Code)

	lreq := httptest.NewRequest(http.MethodPost, "/api/Hosp/staff/login", bytes.NewReader(rb))
	lreq.Header.Set("Content-Type", "application/json")
	lrec := httptest.NewRecorder()
	router.ServeHTTP(lrec, lreq)
	require.Equal(t, http.StatusOK, lrec.Code)
}

func TestStaffRegister_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockAuthService{}
	sh := handlers.NewStaffHandler(mock)
	router := gin.New()
	router.POST("/api/:hospital/staff/create", sh.Register)

	req := httptest.NewRequest(http.MethodPost, "/api/H/staff/create", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestStaffLogin_WrongCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockAuthService{LoginFn: func(hospital, username, password string) (string, *models.Staff, error) {
		return "", nil, errors.New("invalid")
	}}
	sh := handlers.NewStaffHandler(mock)
	router := gin.New()
	router.POST("/api/:hospital/staff/login", sh.Login)

	bad := map[string]string{"username": "x", "password": "y"}
	b, _ := json.Marshal(bad)
	req := httptest.NewRequest(http.MethodPost, "/api/H/staff/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestStaffRegister_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockAuthService{RegisterFn: func(hospital, username, password string) (*models.Staff, error) {
		return nil, errors.New("service fail")
	}}

	sh := handlers.NewStaffHandler(mock)
	router := gin.New()
	router.POST("/api/:hospital/staff/create", sh.Register)

	reg := map[string]string{"username": "u1", "password": "p"}
	rb, _ := json.Marshal(reg)
	rreq := httptest.NewRequest(http.MethodPost, "/api/Hosp/staff/create", bytes.NewReader(rb))
	rreq.Header.Set("Content-Type", "application/json")
	rrec := httptest.NewRecorder()
	router.ServeHTTP(rrec, rreq)
	require.Equal(t, http.StatusInternalServerError, rrec.Code)
}
