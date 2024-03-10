package adapters

import (
	"github.com/google/uuid"
	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"gorm.io/gorm"
)

type UserAdapter struct {
	DB *gorm.DB
}

func NewUserAdapter(db *gorm.DB) *UserAdapter {
	return &UserAdapter{
		DB: db,
	}
}
func (user *UserAdapter) UserSignup(userData entities.User) (entities.User, error) {
	var res entities.User
	id := uuid.New()
	insertQuery := `INSERT INTO users (id,name,email,password,phone) VALUES ($1,$2,$3,$4,$5) RETURNING *`
	if err := user.DB.Raw(insertQuery, id, userData.Name, userData.Email, userData.Password, userData.Phone).Scan(&res).Error; err != nil {
		return entities.User{}, err
	}
	return res, nil
}
func (user *UserAdapter) GetUserByEmail(email string) (entities.User, error) {
	var res entities.User
	selectQuery := `SELECT * FROM USERS WHERE email=?`
	if err := user.DB.Raw(selectQuery, email).Scan(&res).Error; err != nil {
		return entities.User{}, err
	}
	return res, nil
}
func (user *UserAdapter) GetAdminByEmail(email string) (entities.Admin, error) {
	var res entities.Admin
	selectQuery := `SELECT * FROM admins WHERE email=?`
	if err := user.DB.Raw(selectQuery, email).Scan(&res).Error; err != nil {
		return entities.Admin{}, err
	}
	return res, nil
}
