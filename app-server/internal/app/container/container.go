package container

import (
	"github.com/maulana1k/forum-app/gen/recommender"
	"github.com/maulana1k/forum-app/internal/domain/repository"
	"github.com/maulana1k/forum-app/internal/domain/service"
	"github.com/maulana1k/forum-app/internal/provider/broker"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Container struct {
	service.AuthService
	service.UserService
	service.PostService
	service.RecommendationService
}

func NewContainer(db *gorm.DB, grpc *grpc.ClientConn, broker *broker.RabbitMQ) *Container {
	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)

	recClient := recommender.NewRecommenderServiceClient(grpc)

	recRepo := repository.NewRecommendationRepository(recClient)

	return &Container{
		AuthService:           service.NewAuthService(authRepo),
		UserService:           service.NewUserService(userRepo),
		PostService:           service.NewPostService(postRepo, broker),
		RecommendationService: service.NewRecommendationService(recRepo),
	}
}
