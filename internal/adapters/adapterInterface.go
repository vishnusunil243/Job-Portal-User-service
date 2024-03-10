package adapters

import "github.com/vishnusunil243/Job-Portal-User-service/entities"

type AdapterInterface interface {
	UserSignup(entities.User) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
	GetAdminByEmail(email string) (entities.Admin, error)
}
