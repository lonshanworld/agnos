package services

import "agnos_candidate_assignment/models"

// AuthServiceInterface defines the methods used by handlers for authentication
type AuthServiceInterface interface {
	Register(hospital, username, password string) (*models.Staff, error)
	Login(hospital, username, password string) (string, *models.Staff, error)
}

// PatientServiceInterface defines the methods used by handlers for patient operations
type PatientServiceInterface interface {
	Search(hospitalID uint, filters map[string]interface{}) ([]models.Patient, error)
	GetByNationalOrPassport(hospitalID uint, id string) (*models.Patient, error)
}
