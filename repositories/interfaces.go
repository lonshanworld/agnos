package repositories

import "agnos_candidate_assignment/models"

// HospitalRepositoryInterface defines the methods handlers need from hospital repo
type HospitalRepositoryInterface interface {
	Create(h *models.Hospital) error
	FindByName(name string) (*models.Hospital, error)
	FindByID(id uint) (*models.Hospital, error)
}
