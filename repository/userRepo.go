package repository

import (
	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}
func (user *UserRepo) UserSignup(userData entities.User) (entities.User, error) {
	var res entities.User
	insertQuery := `INSERT INTO users (name,email,password,phone) VALUES ($1,$2,$3,$4) RETURNING *`
	if err := user.DB.Raw(insertQuery, userData.Name, userData.Email, userData.Password, userData.Phone).Scan(&res).Error; err != nil {
		return entities.User{}, err
	}
	return res, nil
}
