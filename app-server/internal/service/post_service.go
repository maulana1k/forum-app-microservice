package service

import (
	"errors"
	"time"

	"github.com/maulana1k/forum-app/internal/dto"
	"github.com/maulana1k/forum-app/internal/models"
	"github.com/maulana1k/forum-app/internal/repository"
	"gorm.io/gorm"
)

type PostService interface {
	CreatePost(userID uint, req *dto.CreatePostRequest) (*dto.PostResponse, error)
	GetPostByID(id uint) (*dto.PostResponse, error)
	GetAllPosts(page, limit int) (*dto.PaginatedPostsResponse, error)
	UpdatePost(postID, userID uint, req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	DeletePost(postID, userID uint) error
	GetPostsByUserID(userID uint, page, limit int) (*dto.PaginatedPostsResponse, error)
	LikePost(postID, userID uint) error
	UnlikePost(postID, userID uint) error
	BookmarkPost(postID, userID uint) error
	UnbookmarkPost(postID, userID uint) error
	RepostPost(postID, userID uint) error
	UnrepostPost(postID, userID uint) error
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{
		postRepo: postRepo,
	}
}

func (s *postService) CreatePost(userID uint, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	post := &models.Post{
		Content:  req.Content,
		Tags:     req.Tags,
		ImageURL: req.ImageURL,
		// AuthorID: userID, // Note: Based on your model, this should be a User reference
	}

	if req.QuotedPostID != nil {
		// Verify quoted post exists
		quotedPost, err := s.postRepo.GetPostByID(*req.QuotedPostID)
		if err != nil {
			return nil, errors.New("quoted post not found")
		}
		post.QuotedPostID = &quotedPost.ID
	}

	if err := s.postRepo.CreatePost(post); err != nil {
		return nil, err
	}

	return s.convertPostToResponse(post), nil
}

func (s *postService) GetPostByID(id uint) (*dto.PostResponse, error) {
	post, err := s.postRepo.GetPostWithDetails(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	return s.convertPostToResponse(post), nil
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
		postResponses[i] = *s.convertPostToResponse(&post)
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

func (s *postService) UpdatePost(postID, userID uint, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	// Check if post exists and user is the author
	existingPost, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	// Note: You'll need to implement proper user authorization check
	if existingPost.AuthorID != userID {
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

	return s.convertPostToResponse(updatedPost), nil
}

func (s *postService) DeletePost(postID, userID uint) error {
	// Check if post exists and user is the author
	existingPost, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	// Note: Implement proper user authorization check
	if existingPost.AuthorID != userID {
		return errors.New("unauthorized to delete this post")
	}

	return s.postRepo.DeletePost(postID)
}

func (s *postService) GetPostsByUserID(userID uint, page, limit int) (*dto.PaginatedPostsResponse, error) {
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
		postResponses[i] = *s.convertPostToResponse(&post)
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

func (s *postService) LikePost(postID, userID uint) error {
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

func (s *postService) UnlikePost(postID, userID uint) error {
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

func (s *postService) BookmarkPost(postID, userID uint) error {
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

func (s *postService) UnbookmarkPost(postID, userID uint) error {
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

func (s *postService) RepostPost(postID, userID uint) error {
	// Check if post exists
	_, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return errors.New("post not found")
	}

	// Check if already reposted
	isReposted, err := s.postRepo.IsPostRepostedByUser(postID, userID)
	if err != nil {
		return err
	}
	if isReposted {
		return errors.New("post already reposted")
	}

	return s.postRepo.RepostByUser(postID, userID)
}

func (s *postService) UnrepostPost(postID, userID uint) error {
	// Check if post is reposted
	isReposted, err := s.postRepo.IsPostRepostedByUser(postID, userID)
	if err != nil {
		return err
	}
	if !isReposted {
		return errors.New("post not reposted")
	}

	return s.postRepo.UnrepostByUser(postID, userID)
}

func (s *postService) convertPostToResponse(post *models.Post) *dto.PostResponse {
	response := &dto.PostResponse{
		ID:        post.ID,
		Content:   post.Content,
		Tags:      post.Tags,
		ImageURL:  post.ImageURL,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		// Author:    // You'll need to implement this based on your User model
		Replies: make([]dto.ReplyResponse, len(post.Replies)),
	}

	// Convert replies
	for i, reply := range post.Replies {
		response.Replies[i] = dto.ReplyResponse{
			ID:        reply.ID,
			Content:   reply.Content,
			Author:    reply.Author,
			CreatedAt: reply.CreatedAt,
		}
	}

	// Convert quoted post if exists
	// if post.QuotedPostID != nil {
	// 	response.QuotedPost = &dto.PostResponse{
	// 		ID:        post.QuotedPostID.ID,
	// 		Content:   post.QuotedPostID.Content,
	// 		Tags:      post.QuotedPostID.Tags,
	// 		ImageURL:  post.QuotedPostID.ImageURL,
	// 		CreatedAt: post.QuotedPostID.CreatedAt,
	// 		UpdatedAt: post.QuotedPostID.UpdatedAt,
	// 	}
	// }

	return response
}
