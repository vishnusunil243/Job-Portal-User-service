package initializer

import (
	"github.com/vishnusunil243/Job-Portal-User-service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/service"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB) *service.UserService {
	repo := adapters.NewUserAdapter(db)
	service := service.NewUserService(repo)
	return service
}
