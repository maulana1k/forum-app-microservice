package repository

import (
	"github.com/google/uuid"
	"github.com/maulana1k/forum-app/internal/domain/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	CreatePost(post *models.Post) error
	GetPostByID(id string) (*models.Post, error)
	GetAllPosts(offset, limit int) ([]models.Post, int64, error)
	UpdatePost(id string, post *models.Post) error
	DeletePost(id string) error
	GetPostsByUserID(userID string, offset, limit int) ([]models.Post, int64, error)
	LikePost(postID, userID string) error
	UnlikePost(postID, userID string) error
	IsPostLikedByUser(postID, userID string) (bool, error)
	BookmarkPost(postID, userID string) error
	UnbookmarkPost(postID, userID string) error
	IsPostBookmarkedByUser(postID, userID string) (bool, error)
	// RepostByUser(postID, userID uint) error
	// UnrepostByUser(postID, userID uint) error
	// IsPostRepostedByUser(postID, userID uint) (bool, error)
	GetPostWithDetails(id string) (*models.Post, error)
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

func (r *postRepository) GetPostByID(id string) (*models.Post, error) {
	postID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var post models.Post
	err = r.db.Preload("Author").
		Preload("Replies").
		First(&post, postID).Error
	if err != nil {
		return nil, err
	}
	post.LikesCount = r.getLikesCounts(postID)[postID]
	post.RepliesCount = r.getRepliesCounts(postID)[postID]
	post.RepostsCount = r.getRepostsCounts(postID)[postID]

	return &post, nil
}

func (r *postRepository) GetAllPosts(offset, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// Total posts count
	// if err := r.db.Model(&models.Post{}).Count(&total).Error; err != nil {
	// 	return nil, 0, err
	// }

	// Fetch posts with relationships
	if err := r.db.Preload("Author").
		Preload("Replies").
		Preload("QuotedPost").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	// Extract post IDs
	total = int64(len(posts))
	postIDs := make([]uuid.UUID, len(posts))
	for i, p := range posts {
		postIDs[i] = p.ID
	}

	// Batch fetch counts
	likesCounts := r.getLikesCounts(postIDs...)
	repliesCounts := r.getRepliesCounts(postIDs...)
	repostsCounts := r.getRepostsCounts(postIDs...)

	// Assign counts to posts
	for i, p := range posts {
		posts[i].LikesCount = likesCounts[p.ID]
		posts[i].RepliesCount = repliesCounts[p.ID]
		posts[i].RepostsCount = repostsCounts[p.ID]
	}

	return posts, total, nil
}

// Count likes for a post
func (r *postRepository) getLikesCounts(postIDs ...uuid.UUID) map[uuid.UUID]int {
	if len(postIDs) == 0 {
		return map[uuid.UUID]int{}
	}

	var results []struct {
		PostID uuid.UUID
		Count  int64
	}

	r.db.Model(&models.PostInteractions{}).
		Select("post_id, COUNT(*) as count").
		Where("post_id IN ?", postIDs).
		Group("post_id").
		Scan(&results)

	counts := make(map[uuid.UUID]int, len(results))
	for _, r := range results {
		counts[r.PostID] = int(r.Count)
	}
	return counts
}

// Batch count replies
func (r *postRepository) getRepliesCounts(postIDs ...uuid.UUID) map[uuid.UUID]int {
	if len(postIDs) == 0 {
		return map[uuid.UUID]int{}
	}

	var results []struct {
		PostID uuid.UUID
		Count  int64
	}
	r.db.Model(&models.Replies{}).
		Select("post_id, COUNT(*) as count").
		Where("post_id IN ?", postIDs).
		Group("post_id").
		Scan(&results)

	counts := make(map[uuid.UUID]int)
	for _, r := range results {
		counts[r.PostID] = int(r.Count)
	}
	return counts
}

// Batch count reposts
func (r *postRepository) getRepostsCounts(postIDs ...uuid.UUID) map[uuid.UUID]int {
	if len(postIDs) == 0 {
		return map[uuid.UUID]int{}
	}
	var results []struct {
		PostID uuid.UUID
		Count  int64
	}
	r.db.Model(&models.Post{}).
		Select("quoted_post_id as post_id, COUNT(*) as count").
		Where("quoted_post_id IN ?", postIDs).
		Group("quoted_post_id").
		Scan(&results)

	counts := make(map[uuid.UUID]int)
	for _, r := range results {
		counts[r.PostID] = int(r.Count)
	}
	return counts
}

func (r *postRepository) UpdatePost(id string, post *models.Post) error {
	return r.db.Model(&models.Post{}).Where("id = ?", id).Updates(post).Error
}

func (r *postRepository) DeletePost(id string) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *postRepository) GetPostsByUserID(userID string, offset, limit int) ([]models.Post, int64, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, err
	}

	var posts []models.Post
	var total int64

	if err := r.db.Model(&models.Post{}).
		Where("author_id = ?", uid).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Where("author_id = ?", uid).
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "display_name", "avatar_url", "bio")
		}).
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("QuotedPost").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) LikePost(postIDstr, userIDstr string) error {
	postID, err := uuid.Parse(postIDstr)
	if err != nil {
		return err
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return err
	}
	like := models.PostInteractions{
		PostID:          postID,
		UserID:          userID,
		InteractionType: string(models.LIKE),
	}
	return r.db.Create(&like).Error
}

func (r *postRepository) UnlikePost(postIDstr, userIDstr string) error {
	postID, err := uuid.Parse(postIDstr)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return err
	}
	return r.db.Where("post_id = ? AND user_id = ? AND interaction_type = 'LIKE'", postID, userID).
		Delete(&models.PostInteractions{}).Error
}

func (r *postRepository) IsPostLikedByUser(postID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.PostInteractions{}).
		Where("post_id = ? AND user_id = ? AND interaction_type = 'LIKE'", postID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *postRepository) BookmarkPost(postIDstr, userIDstr string) error {
	postID, err := uuid.Parse(postIDstr)
	if err != nil {
		return err
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return err
	}
	bookmark := models.PostInteractions{
		PostID:          postID,
		UserID:          userID,
		InteractionType: string(models.BOOKMARK),
	}
	return r.db.Create(&bookmark).Error
}

func (r *postRepository) UnbookmarkPost(postIDstr, userIDstr string) error {
	postID, err := uuid.Parse(postIDstr)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return err
	}
	return r.db.Where("post_id = ? AND user_id = ? AND interaction_type = 'BOOKMARK'", postID, userID).
		Delete(&models.PostInteractions{}).Error
}

func (r *postRepository) IsPostBookmarkedByUser(postID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.PostInteractions{}).
		Where("post_id = ? AND user_id = ? AND interaction_type = 'BOOKMARK'", postID, userID).
		Count(&count).Error
	return count > 0, err
}

// func (r *postRepository) RepostByUser(postID, userID uint) error {
// 	repost := models.RepostByUser{
// 		PostID: postID,
// 		UserID: userID,
// 	}
// 	return r.db.Create(&repost).Error
// }

// func (r *postRepository) UnrepostByUser(postID, userID uint) error {
// 	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.RepostByUser{}).Error
// }

// func (r *postRepository) IsPostRepostedByUser(postID, userID uint) (bool, error) {
// 	var count int64
// 	err := r.db.Model(&models.RepostByUser{}).
// 		Where("post_id = ? AND user_id = ?", postID, userID).
// 		Count(&count).Error
// 	return count > 0, err
// }

func (r *postRepository) GetPostWithDetails(id string) (*models.Post, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var post models.Post
	err = r.db.Preload("Author").
		Preload("Replies.Post").
		Preload("QuotedPost.Author").
		First(&post, uid).Error
	if err != nil {
		return nil, err
	}

	post.LikesCount = r.getLikesCounts(uid)[uid]
	post.RepliesCount = r.getRepliesCounts(uid)[uid]
	post.RepostsCount = r.getRepostsCounts(uid)[uid]

	return &post, nil
}
