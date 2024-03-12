package adapters

import (
	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"github.com/vishnusunil243/Job-Portal-User-service/entities/helperstruct"
)

type AdapterInterface interface {
	UserSignup(entities.User) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
	GetAdminByEmail(email string) (entities.Admin, error)
	CreateProfile(userID string) error
	GetProfileIdByUserId(userId string) (string, error)
	AdminAddSkill(entities.Skill) error
	UserAddSkill(entities.UserSkill) error
	AddLink(entities.Link) error
	AdminAddCategory(entities.Category) error
	AdminUpdateCategory(entities.Category) error
	GetAllCategory() ([]entities.Category, error)
	AdminUpdateSkill(entities.Skill) error
	AdminGetAllSkills() ([]helperstruct.SkillHelper, error)
	UserDeleteSkill(entities.UserSkill) error
	DeleteLink(string) error
	UserGetAllSkills(profileId string) ([]helperstruct.SkillHelper, error)
	GetAllLinksUser(profileID string) ([]entities.Link, error)
}
