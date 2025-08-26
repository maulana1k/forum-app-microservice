package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID           uint   `gorm:"primaryKey"`
	Content      string `gorm:"type:text;not null"`
	AuthorID     uint   `gorm:"not null;index"`
	QuotedPostID *uint
	Tags         string `gorm:"type:text"`
	ImageURL     string `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	Author     User  `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	QuotedPost *Post `gorm:"foreignKey:QuotedPostID;constraint:OnDelete:SET NULL"`

	Replies []Replies `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Replies struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint `gorm:"not null;index"`
	ParentID  *uint
	Author    string `gorm:"not null"`
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Post   Post     `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Parent *Replies `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
}

type PostLikes struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint `gorm:"not null;index"`
	UserID    uint `gorm:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type RepostByUser struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint `gorm:"not null;index"`
	UserID    uint `gorm:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type BookmarkByUser struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint `gorm:"not null;index"`
	UserID    uint `gorm:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
