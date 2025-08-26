package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	AvatarURL string `gorm:"type:text"`
	Role      string `gorm:"default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Session struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	Token     string `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type UserProfile struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;uniqueIndex"`
	Bio       string `gorm:"type:text"`
	Location  string `gorm:"type:varchar(100)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
type UserSettings struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;uniqueIndex"`
	Theme     string `gorm:"type:varchar(50);default:'light'"`
	Language  string `gorm:"type:varchar(50);default:'en'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
