package services

import (
	"agnos_candidate_assignment/config"
	"agnos_candidate_assignment/models"
	"agnos_candidate_assignment/repositories"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	StaffRepo    *repositories.StaffRepository
	HospitalRepo *repositories.HospitalRepository
	conf         *config.Config
}

func NewAuthService(staffRepo *repositories.StaffRepository, hospitalRepo *repositories.HospitalRepository, conf *config.Config) *AuthService {
	return &AuthService{
		StaffRepo:    staffRepo,
		HospitalRepo: hospitalRepo,
		conf:         conf,
	}
}

func (auth *AuthService) HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (auth *AuthService) CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (auth *AuthService) Register(hospitalName, userName, password string) (*models.Staff, error) {
	hospital, err := auth.HospitalRepo.FindByName(hospitalName)
	if err != nil {
		return nil, errors.New("Hospital not found")
	}

	if _, err := auth.StaffRepo.GetByUsenameAndHospital(userName, hospital.ID); err == nil {
		return nil, errors.New("Username already exists in this hospital")
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	staff := &models.Staff{UserName: userName, PasswordHash: hashedPassword, HospitalID: hospital.ID}
	if err := auth.StaffRepo.CreateStaff(staff); err != nil {
		return nil, err
	}
	return staff, nil
}

type StaffClaims struct {
	StaffID    uint
	HospitalID uint
	jwt.RegisteredClaims
}

func (auth *AuthService) Login(hospitalName, username, password string) (string, *models.Staff, error) {
	hospital, err := auth.HospitalRepo.FindByName(hospitalName)
	if err != nil {
		return "", nil, errors.New("Hospital not found")
	}

	staff, err := auth.StaffRepo.GetByUsenameAndHospital(username, hospital.ID)
	if err != nil {
		return "", nil, errors.New("Invalid username or password")
	}

	if err := auth.CheckPasswordHash(password, staff.PasswordHash); err != nil {
		return "", nil, errors.New("Invalid username or password")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"staff_id":    staff.ID,
		"hospital_id": staff.HospitalID,
		"iat":         now.Unix(),
		"exp":         now.Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(auth.conf.JwtSecret))
	if err != nil {
		return "", nil, err
	}

	return signedToken, staff, nil
}
