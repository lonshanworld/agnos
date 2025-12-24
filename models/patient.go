package models

import "time"

type Patient struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	HospitalID   uint      `gorm:"not null;index" json:"hospital_id"`
	Hospital     Hospital  `gorm:"foreignKey:HospitalID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"hospital,omitempty"`
	FirstNameTH  string    `gorm:"size:255;not null" json:"first_name_th"`
	MiddleNameTH string    `gorm:"size:255" json:"middle_name_th"`
	LastNameTH   string    `gorm:"size:255" json:"last_name_th"`
	DateOfBirth  time.Time `gorm:"type:date;not null" json:"date_of_birth"`
	PatientHN    string     `gorm:"size:50;uniqueIndex" json:"patient_hn"`
	NationalID   *string    `gorm:"size:255;uniqueIndex" json:"national_id,omitempty"`
	PassportID   *string    `gorm:"size:255;uniqueIndex" json:"passport_id,omitempty"`
	PhoneNumber  *string    `gorm:"size:50" json:"phone_number,omitempty"`
	Email        *string    `gorm:"size:255" json:"email,omitempty"`
	Gender       Gender    `gorm:"size:1;not null" json:"gender"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
