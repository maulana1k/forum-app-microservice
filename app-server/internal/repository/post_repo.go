package repository

import (
	"github.com/maulana1k/forum-app/internal/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	CreatePost(post *models.Post) error
	GetPostByID(id uint) (*models.Post, error)
	GetAllPosts(offset, limit int) ([]models.Post, int64, error)
	UpdatePost(id uint, post *models.Post) error
	DeletePost(id uint) error
	GetPostsByUserID(userID uint, offset, limit int) ([]models.Post, int64, error)
	LikePost(postID, userID uint) error
	UnlikePost(postID, userID uint) error
	IsPostLikedByUser(postID, userID uint) (bool, error)
	BookmarkPost(postID, userID uint) error
	UnbookmarkPost(postID, userID uint) error
	IsPostBookmarkedByUser(postID, userID uint) (bool, error)
	RepostByUser(postID, userID uint) error
	UnrepostByUser(postID, userID uint) error
	IsPostRepostedByUser(postID, userID uint) (bool, error)
	GetPostWithDetails(id uint) (*models.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) CreatePost(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("AuthorID").Preload("Replies").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetAllPosts(offset, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// Get total count
	if err := r.db.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get posts with pagination and preloaded relationships
	err := r.db.Preload("AuthorID").
		Preload("Replies").
		Preload("QuotedPostID").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) UpdatePost(id uint, post *models.Post) error {
	return r.db.Model(&models.Post{}).Where("id = ?", id).Updates(post).Error
}

func (r *postRepository) DeletePost(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *postRepository) GetPostsByUserID(userID uint, offset, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// Get total count for user posts
	if err := r.db.Model(&models.Post{}).Where("author_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get user posts with pagination
	err := r.db.Preload("AuthorID").
		Preload("Replies").
		Preload("QuotedPostID").
		Where("author_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) LikePost(postID, userID uint) error {
	like := models.PostLikes{
		PostID: postID,
		UserID: userID,
	}
	return r.db.Create(&like).Error
}

func (r *postRepository) UnlikePost(postID, userID uint) error {
	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.PostLikes{}).Error
}

func (r *postRepository) IsPostLikedByUser(postID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.PostLikes{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *postRepository) BookmarkPost(postID, userID uint) error {
	bookmark := models.BookmarkByUser{
		PostID: postID,
		UserID: userID,
	}
	return r.db.Create(&bookmark).Error
}

func (r *postRepository) UnbookmarkPost(postID, userID uint) error {
	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.BookmarkByUser{}).Error
}

func (r *postRepository) IsPostBookmarkedByUser(postID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.BookmarkByUser{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *postRepository) RepostByUser(postID, userID uint) error {
	repost := models.RepostByUser{
		PostID: postID,
		UserID: userID,
	}
	return r.db.Create(&repost).Error
}

func (r *postRepository) UnrepostByUser(postID, userID uint) error {
	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.RepostByUser{}).Error
}

func (r *postRepository) IsPostRepostedByUser(postID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.RepostByUser{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *postRepository) GetPostWithDetails(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("AuthorID").
		Preload("Replies.Post").
		Preload("QuotedPostID.AuthorID").
		First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}
