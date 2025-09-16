package repository

import (
	"context"

	pb "github.com/maulana1k/forum-app/gen/recommender"
)

type RecommendationRepository interface {
	GetRecommendedPosts(userID, topic string, limit int) ([]*pb.PostItem, error)
}

type recommendationRepository struct {
	client pb.RecommenderServiceClient
}

func NewRecommendationRepository(client pb.RecommenderServiceClient) RecommendationRepository {
	return &recommendationRepository{client: client}
}

func (r *recommendationRepository) GetRecommendedPosts(userID, topic string, limit int) ([]*pb.PostItem, error) {
	req := &pb.RecommendationRequest{
		UserId: userID,
		Topic:  topic,
		Limit:  int32(limit),
	}

	resp, err := r.client.GetRecommendedPosts(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp.Posts, nil
}
