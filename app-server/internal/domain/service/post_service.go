package service

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/maulana1k/forum-app/internal/app/dto"
	"github.com/maulana1k/forum-app/internal/domain/models"
	"github.com/maulana1k/forum-app/internal/domain/repository"
	"github.com/maulana1k/forum-app/internal/provider/broker"
	"gorm.io/gorm"
)

type PostService interface {
	CreatePost(userID string, req *dto.CreatePostRequest) (*dto.PostResponse, error)
	GetPostByID(id string) (*dto.PostResponse, error)
	GetAllPosts(page, limit int) (*dto.PaginatedPostsResponse, error)
	UpdatePost(postID, userID string, req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	DeletePost(postID, userID string) error
	GetPostsByUserID(userID string, page, limit int) (*dto.PaginatedPostsResponse, error)
	LikePost(postID, userID string) error
	UnlikePost(postID, userID string) error
	BookmarkPost(postID, userID string) error
	UnbookmarkPost(postID, userID string) error
}

type postService struct {
	postRepo repository.PostRepository
	broker   *broker.RabbitMQ
}

func NewPostService(postRepo repository.PostRepository, brokerc *broker.RabbitMQ) PostService {
	return &postService{
		postRepo: postRepo,
		broker:   brokerc,
	}
}

func (s *postService) CreatePost(userID string, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	id, _ := uuid.Parse(userID)
	post := &models.Post{
		Content:  req.Content,
		Tags:     req.Tags,
		ImageURL: req.ImageURL,
		AuthorID: id,
	}

	if req.QuotedPostID != "" {
		// Verify quoted post exists
		quotedPost, err := s.postRepo.GetPostByID(req.QuotedPostID)
		if err != nil {
			return nil, errors.New("quoted post not found")
		}
		post.QuotedPostID = &quotedPost.ID
	}

	if err := s.postRepo.CreatePost(post); err != nil {
		return nil, err
	}

	event := map[string]any{
		"post_id": post.ID,
		"author":  post.AuthorID,
		"content": post.Content,
	}
	body, _ := json.Marshal(event)

	producer := broker.NewProducer(s.broker, "post-create")

	err := producer.Publish(body)
	if err != nil {
		log.Printf("Failed to publish post message: %v", err)
	}

	return s.MapPostToResponse(post), nil
}

func (s *postService) GetPostByID(id string) (*dto.PostResponse, error) {
	post, err := s.postRepo.GetPostWithDetails(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	return s.MapPostToResponse(post), nil
}

func (s *postService) GetAllPosts(page, limit int) (*dto.PaginatedPostsResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	posts, total, err := s.postRepo.GetAllPosts(offset, limit)
	if err != nil {
		return nil, err
	}

	postResponses := make([]dto.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = *s.MapPostToResponse(&post)
	}

	totalPages := (int(total) + limit - 1) / limit

	return &dto.PaginatedPostsResponse{
		Posts:       postResponses,
		Total:       int(total),
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
	}, nil
}

func (s *postService) UpdatePost(postID, userID string, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	// Check if post exists and user is the author
	existingPost, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	// Note: You'll need to implement proper user authorization check
	if existingPost.AuthorID.String() != userID {
		return nil, errors.New("unauthorized to update this post")
	}

	updateData := &models.Post{}
	if req.Content != nil {
		updateData.Content = *req.Content
	}
	if req.Tags != nil {
		updateData.Tags = *req.Tags
	}
	if req.ImageURL != nil {
		updateData.ImageURL = *req.ImageURL
	}
	updateData.UpdatedAt = time.Now()

	if err := s.postRepo.UpdatePost(postID, updateData); err != nil {
		return nil, err
	}

	updatedPost, err := s.postRepo.GetPostWithDetails(postID)
	if err != nil {
		return nil, err
	}

	return s.MapPostToResponse(updatedPost), nil
}

func (s *postService) DeletePost(postID, userID string) error {
	// Check if post exists and user is the author
	existingPost, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	// Note: Implement proper user authorization check
	if existingPost.AuthorID.String() != userID {
		return errors.New("unauthorized to delete this post")
	}

	return s.postRepo.DeletePost(postID)
}

func (s *postService) GetPostsByUserID(userID string, page, limit int) (*dto.PaginatedPostsResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	posts, total, err := s.postRepo.GetPostsByUserID(userID, offset, limit)
	if err != nil {
		return nil, err
	}

	postResponses := make([]dto.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = *s.MapPostToResponse(&post)
	}

	totalPages := (int(total) + limit - 1) / limit

	return &dto.PaginatedPostsResponse{
		Posts:       postResponses,
		Total:       int(total),
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
	}, nil
}

func (s *postService) LikePost(postID, userID string) error {
	// Check if post exists
	_, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return errors.New("post not found")
	}

	// Check if already liked
	isLiked, err := s.postRepo.IsPostLikedByUser(postID, userID)
	if err != nil {
		return err
	}
	if isLiked {
		return errors.New("post already liked")
	}

	return s.postRepo.LikePost(postID, userID)
}

func (s *postService) UnlikePost(postID, userID string) error {
	// Check if post is liked
	isLiked, err := s.postRepo.IsPostLikedByUser(postID, userID)
	if err != nil {
		return err
	}
	if !isLiked {
		return errors.New("post not liked")
	}

	return s.postRepo.UnlikePost(postID, userID)
}

func (s *postService) BookmarkPost(postID, userID string) error {
	// Check if post exists
	_, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return errors.New("post not found")
	}

	// Check if already bookmarked
	isBookmarked, err := s.postRepo.IsPostBookmarkedByUser(postID, userID)
	if err != nil {
		return err
	}
	if isBookmarked {
		return errors.New("post already bookmarked")
	}

	return s.postRepo.BookmarkPost(postID, userID)
}

func (s *postService) UnbookmarkPost(postID, userID string) error {
	// Check if post is bookmarked
	isBookmarked, err := s.postRepo.IsPostBookmarkedByUser(postID, userID)
	if err != nil {
		return err
	}
	if !isBookmarked {
		return errors.New("post not bookmarked")
	}

	return s.postRepo.UnbookmarkPost(postID, userID)
}

// func (s *postService) RepostPost(postID, userID uint) error {
// 	// Check if post exists
// 	_, err := s.postRepo.GetPostByID(postID)
// 	if err != nil {
// 		return errors.New("post not found")
// 	}

// 	// Check if already reposted
// 	isReposted, err := s.postRepo.IsPostRepostedByUser(postID, userID)
// 	if err != nil {
// 		return err
// 	}
// 	if isReposted {
// 		return errors.New("post already reposted")
// 	}

// 	return s.postRepo.RepostByUser(postID, userID)
// }

// func (s *postService) UnrepostPost(postID, userID uint) error {
// 	// Check if post is reposted
// 	isReposted, err := s.postRepo.IsPostRepostedByUser(postID, userID)
// 	if err != nil {
// 		return err
// 	}
// 	if !isReposted {
// 		return errors.New("post not reposted")
// 	}

// 	return s.postRepo.UnrepostByUser(postID, userID)
// }

func (s *postService) MapPostToResponse(p *models.Post) *dto.PostResponse {
	author := dto.PostAuthor{
		ID:          p.Author.ID.String(),
		Username:    p.Author.Username,
		DisplayName: p.Author.DisplayName,
		AvatarURL:   p.Author.AvatarURL,
		Bio:         p.Author.Bio,
	}

	replies := make([]dto.ReplyResponse, len(p.Replies))
	for i, r := range p.Replies {
		replies[i] = dto.ReplyResponse{
			ID:        r.ID,
			Content:   r.Content,
			Author:    r.Author,
			CreatedAt: r.CreatedAt,
		}
	}

	quotedPostID := ""
	if p.QuotedPost != nil {
		quotedPostID = p.QuotedPost.ID.String()
	}

	return &dto.PostResponse{
		ID:           p.ID.String(),
		Content:      p.Content,
		Tags:         p.Tags,
		ImageURL:     p.ImageURL,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
		Author:       author,
		Replies:      replies,
		QuotedPost:   quotedPostID,
		LikesCount:   p.LikesCount,
		RepliesCount: p.RepliesCount,
		RepostsCount: p.RepostsCount,
		// set these flags according to your business logic
		IsLiked:      false,
		IsBookmarked: false,
		IsReposted:   false,
	}
}
