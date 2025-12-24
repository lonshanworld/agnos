package repositories

import (
	"agnos_candidate_assignment/models"

	"gorm.io/gorm"
)

type HospitalRepository struct {
	db *gorm.DB
}

func NewHospitalRepository(db *gorm.DB) *HospitalRepository {
	return &HospitalRepository{db: db}
}

func (r *HospitalRepository) Create(h *models.Hospital) error {
	return r.db.Create(h).Error
}

func (r *HospitalRepository) FindByName(name string) (*models.Hospital, error) {
	var h models.Hospital
	if err := r.db.Where("name = ?", name).First(&h).Error; err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *HospitalRepository) FindByID(id uint) (*models.Hospital, error) {
	var h models.Hospital
	if err := r.db.First(&h, id).Error; err != nil {
		return nil, err
	}
	return &h, nil
}
