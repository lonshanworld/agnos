package repositories

import "agnos_candidate_assignment/models"

type HospitalRepositoryInterface interface {
	Create(h *models.Hospital) error
	FindByName(name string) (*models.Hospital, error)
	FindByID(id uint) (*models.Hospital, error)
}
