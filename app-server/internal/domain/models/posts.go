package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Content      string    `gorm:"type:text;not null"`
	AuthorID     uuid.UUID `gorm:"not null;index"`
	QuotedPostID *uuid.UUID
	Tags         string `gorm:"type:text"`
	ImageURL     string `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	LikesCount   int
	RepliesCount int
	RepostsCount int
	Author       User  `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	QuotedPost   *Post `gorm:"foreignKey:QuotedPostID;constraint:OnDelete:SET NULL"`

	Replies []Replies `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Replies struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint `gorm:"not null;index"`
	ParentID  *uint
	Author    string        `gorm:"not null"`
	Content   string        `gorm:"type:text;not null"`
	Likes     pq.Int64Array `gorm:"type:integer[]"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Post   Post     `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Parent *Replies `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
}

type PostInteractionType string

const (
	LIKE     PostInteractionType = "LIKE"
	DISLIKE  PostInteractionType = "DISLIKE"
	BOOKMARK PostInteractionType = "BOOKMARK"
)

type PostInteractions struct {
	ID              uint      `gorm:"primaryKey"`
	PostID          uuid.UUID `gorm:"not null;index"`
	UserID          uuid.UUID `gorm:"not null;index"`
	InteractionType string    `gorm:"not null;type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// type PostLikes struct {
// 	ID        uint `gorm:"primaryKey"`
// 	PostID    uint `gorm:"not null;index"`
// 	UserID    uint `gorm:"not null;index"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`

// 	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
// 	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
// }

// type RepostByUser struct {
// 	ID        uint `gorm:"primaryKey"`
// 	PostID    uint `gorm:"not null;index"`
// 	UserID    uint `gorm:"not null;index"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`

// 	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
// 	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
// }

// type BookmarkByUser struct {
// 	ID        uint `gorm:"primaryKey"`
// 	PostID    uint `gorm:"not null;index"`
// 	UserID    uint `gorm:"not null;index"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`

// 	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
// 	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
// }
