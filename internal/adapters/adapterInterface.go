package adapters

import (
	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"github.com/vishnusunil243/Job-Portal-User-service/entities/helperstruct"
)

type AdapterInterface interface {
	UserSignup(entities.User) (entities.User, error)
	GetUserByPhone(phone string) (entities.User, error)
	GetCategoryByName(category string) (entities.Category, error)
	GetUserByEmail(email string) (entities.User, error)
	GetAdminByEmail(email string) (entities.Admin, error)
	CreateProfile(userID string) error
	GetProfileIdByUserId(userId string) (string, error)
	AdminAddSkill(entities.Skill) error
	UserAddSkill(entities.UserSkill) error
	AddLink(entities.Link) error
	GetLinkByTitle(profileId string, title string) (entities.Link, error)
	GetSkillByName(skill string) (entities.Skill, error)
	AdminAddCategory(entities.Category) error
	AdminUpdateCategory(entities.Category) error
	GetAllCategory() ([]entities.Category, error)
	AdminUpdateSkill(entities.Skill) error
	AdminGetAllSkills() ([]helperstruct.SkillHelper, error)
	UserDeleteSkill(entities.UserSkill) error
	DeleteLink(string) error
	UserGetAllSkills(profileId string) ([]helperstruct.SkillHelper, error)
	GetAllLinksUser(profileID string) ([]entities.Link, error)
	GetSkillById(id int) (helperstruct.SkillHelper, error)
	GetCategoryById(id int) (entities.Category, error)
	GetUserById(userId string) (entities.User, error)
	JobApply(entities.Jobs) error
	UserEditName(entities.User) error
	UserEditPhone(entities.User) error
	UserAddAddress(entities.Address) error
	UserEditAddress(entities.Address) error
	GetAddressByProfileId(profileId string) (entities.Address, error)
	UploadProfileImage(Image, ProfileId string) (string, error)
	GetProfilePic(string) (string, error)
	GetAppliedJobIds(string) ([]string, error)
	GetAppliedJob(string, string) (entities.Jobs, error)
	GetUserSkillById(string, int) (entities.UserSkill, error)
}
