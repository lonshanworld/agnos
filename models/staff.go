package models

import "time"

type Staff struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName     string    `gorm:"size:255;not null;uniqueIndex" json:"user_name"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	HospitalID   uint      `gorm:"not null;index" json:"hospital_id"`
	Hospital     Hospital  `gorm:"foreignKey:HospitalID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"hospital,omitempty"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
