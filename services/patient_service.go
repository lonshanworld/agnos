package services

import (
	"agnos_candidate_assignment/models"
	"agnos_candidate_assignment/repositories"
)

type PatientService struct {
	Repo *repositories.PatientRepository
}

func NewPatientService(repo *repositories.PatientRepository) *PatientService {
	return &PatientService{Repo: repo}
}

func (patientservice *PatientService) Search(hospitalID uint, filters map[string]interface{}) ([]models.Patient, error) {
	return patientservice.Repo.Search(hospitalID, filters)
}

func (patientservice *PatientService) GetByNationalOrPassport(hospitalID uint, nationalOrPassport string) (*models.Patient, error) {
	return patientservice.Repo.GetByNationalOrPassportID(hospitalID, nationalOrPassport)
}
