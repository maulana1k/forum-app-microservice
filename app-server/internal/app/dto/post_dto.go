package dto

import (
	"time"
)

// CreatePostRequest represents the request body for creating a post
type CreatePostRequest struct {
	Content      string `json:"content" validate:"required,min=1,max=2000"`
	Tags         string `json:"tags,omitempty"`
	ImageURL     string `json:"image_url,omitempty" validate:"omitempty,url"`
	QuotedPostID string `json:"quoted_post_id,omitempty"`
}

// UpdatePostRequest represents the request body for updating a post
type UpdatePostRequest struct {
	Content  *string `json:"content,omitempty" validate:"omitempty,min=1,max=2000"`
	Tags     *string `json:"tags,omitempty"`
	ImageURL *string `json:"image_url,omitempty" validate:"omitempty,url"`
}

type PostAuthor struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
}

type PostResponse struct {
	ID           string          `json:"id"`
	Content      string          `json:"content"`
	Tags         string          `json:"tags"`
	ImageURL     string          `json:"image_url"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Author       PostAuthor      `json:"author"`
	Replies      []ReplyResponse `json:"replies"`
	QuotedPost   string          `json:"quoted_post"`
	LikesCount   int             `json:"likes_count"`
	RepliesCount int             `json:"replies_count"`
	RepostsCount int             `json:"reposts_count"`
	IsLiked      bool            `json:"is_liked"`
	IsBookmarked bool            `json:"is_bookmarked"`
	IsReposted   bool            `json:"is_reposted"`
}

// ReplyResponse represents a reply in API responses
type ReplyResponse struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaginatedPostsResponse represents paginated posts response
type PaginatedPostsResponse struct {
	Posts       []PostResponse `json:"posts"`
	Total       int            `json:"total"`
	Page        int            `json:"page"`
	Limit       int            `json:"limit"`
	TotalPages  int            `json:"total_pages"`
	HasNextPage bool           `json:"has_next_page"`
	HasPrevPage bool           `json:"has_prev_page"`
}

// PostActionRequest represents like/bookmark/repost actions
type PostActionRequest struct {
	PostID uint `json:"post_id" validate:"required,min=1"`
}

// PostQueryParams represents query parameters for posts
type PostQueryParams struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	UserID *uint  `query:"user_id" validate:"omitempty,min=1"`
	Tags   string `query:"tags"`
	Search string `query:"search"`
}
