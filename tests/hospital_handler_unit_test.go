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

type mockHospitalRepo struct {
	CreateFn func(h *models.Hospital) error
}

func (m *mockHospitalRepo) Create(h *models.Hospital) error {
	return m.CreateFn(h)
}

func (m *mockHospitalRepo) FindByName(name string) (*models.Hospital, error) {
	return nil, nil
}

func (m *mockHospitalRepo) FindByID(id uint) (*models.Hospital, error) {
	return nil, nil
}

func TestHospitalCreate_Positive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockHospitalRepo{CreateFn: func(h *models.Hospital) error {
		h.ID = 11
		return nil
	}}

	hh := handlers.NewHospitalHandler(mock)

	router := gin.New()
	router.POST("/api/hospital", hh.Create)

	body := map[string]string{"name": "UT Hospital"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/hospital", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	var resp models.Hospital
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, uint(11), resp.ID)
	require.Equal(t, "UT Hospital", resp.Name)
}

func TestHospitalCreate_NegativeCreateError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockHospitalRepo{CreateFn: func(h *models.Hospital) error {
		return errors.New("boom")
	}}

	hh := handlers.NewHospitalHandler(&struct{ *mockHospitalRepo }{mock})
	router := gin.New()
	router.POST("/api/hospital", hh.Create)

	body := map[string]string{"name": "UT Hospital"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/hospital", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
