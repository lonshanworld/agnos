package repositories

import (
	"agnos_candidate_assignment/models"

	"gorm.io/gorm"
)

type StaffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) *StaffRepository {
	return &StaffRepository{db: db}
}

func (repo *StaffRepository) CreateStaff(staff *models.Staff) error {
	return repo.db.Create(staff).Error
}

func (repo *StaffRepository) GetByUsenameAndHospital(username string, hospitalID uint) (*models.Staff, error) {
	var staff models.Staff
	if err := repo.db.Where("user_name = ? AND hospital_id = ?", username, hospitalID).First(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func (repo *StaffRepository) GetByID(id uint) (*models.Staff, error) {
	var staff models.Staff
	if err := repo.db.First(&staff, id).Error; err != nil {
		return nil, err
	}

	return &staff, nil
}
