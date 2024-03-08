package Interface

import "github.com/vishnusunil243/Job-Portal-User-service/entities"

type UserRepo interface {
	UserSignup(entities.User) (entities.User, error)
}
