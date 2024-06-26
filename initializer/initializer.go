package initializer

import (
	"github.com/vishnusunil243/Job-Portal-User-service/concurrency"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/service"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/usecases"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB) *service.UserService {
	repo := adapters.NewUserAdapter(db)
	usecase := usecases.NewUserUseCase(repo)
	service := service.NewUserService(repo, usecase, "notification-service:8087", "company-service:8082")
	c := concurrency.NewConcurrency(db, repo, service)
	c.Concurrency()
	return service
}
