package initializer

import (
	"github.com/vishnusunil243/Job-Portal-User-service/repository"
	"github.com/vishnusunil243/Job-Portal-User-service/service"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB) *service.UserService {
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	return service
}
