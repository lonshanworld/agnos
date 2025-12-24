package repositories

import (
	"agnos_candidate_assignment/models"

	"gorm.io/gorm"
)

type PatientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (repo *PatientRepository) Create(p *models.Patient) error {
	return repo.db.Create(p).Error
}

func (repo *PatientRepository) Search(hospitalID uint, filters map[string]interface{}) ([]models.Patient, error) {
	db := repo.db.Model(&models.Patient{}).Where("hospital_id = ?", hospitalID)
	for k, v := range filters {
		db = db.Where(k+" = ?", v)
	}

	var results []models.Patient
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *PatientRepository) GetByNationalOrPassportID(hospitalID uint, id string) (*models.Patient, error) {
	var result models.Patient
	if err := repo.db.Where("hospital_id = ? AND (national_id = ? OR passport_id = ?)", hospitalID, id, id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
