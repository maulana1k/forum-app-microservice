package service

import (
	"github.com/maulana1k/forum-app/internal/app/dto"
	"github.com/maulana1k/forum-app/internal/domain/repository"
)

type RecommendationService interface {
	GetRecommendedPosts(userID, topic string, limit int) ([]dto.PostResponse, error)
}

type recommendationService struct {
	repo repository.RecommendationRepository
}

func NewRecommendationService(repo repository.RecommendationRepository) RecommendationService {
	return &recommendationService{repo: repo}
}

func (s *recommendationService) GetRecommendedPosts(userID, topic string, limit int) ([]dto.PostResponse, error) {
	posts, err := s.repo.GetRecommendedPosts(userID, topic, limit)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	postResponses := make([]dto.PostResponse, len(posts))
	for i, p := range posts {
		postResponses[i] = dto.PostResponse{
			ID:           p.PostId,
			Content:      p.Content,
			LikesCount:   0,
			RepliesCount: 0,
			RepostsCount: 0,
		}
	}

	return postResponses, nil
}
