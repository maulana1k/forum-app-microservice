package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username    string    `gorm:"uniqueIndex;not null"`
	Email       string    `gorm:"uniqueIndex;not null"`
	Password    string    `gorm:"not null"`
	Role        string    `gorm:"default:'user'"`
	AvatarURL   string    `gorm:"type:text"`
	DisplayName string    `gorm:"type:text"`
	Bio         string    `gorm:"type:text"`
	Location    string    `gorm:"type:varchar(100)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type UserSettings struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"not null;uniqueIndex"`
	Theme     string    `gorm:"type:varchar(50);default:'light'"`
	Language  string    `gorm:"type:varchar(50);default:'en'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
