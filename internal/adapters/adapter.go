package adapters

import (
	"time"

	"github.com/google/uuid"
	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"github.com/vishnusunil243/Job-Portal-User-service/entities/helperstruct"
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
func (user *UserAdapter) GetUserByPhone(phone string) (entities.User, error) {
	var res entities.User
	selectQuery := `SELECT * FROM USERS WHERE phone=?`
	if err := user.DB.Raw(selectQuery, phone).Scan(&res).Error; err != nil {
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
func (user *UserAdapter) CreateProfile(userID string) error {

	profileId := uuid.New()
	insertProfileQuery := `INSERT INTO profiles (id,user_id) VALUES ($1,$2)`
	if err := user.DB.Exec(insertProfileQuery, profileId, userID).Error; err != nil {
		return err
	}
	return nil

}
func (user *UserAdapter) GetProfileIdByUserId(userId string) (string, error) {
	var profileId string
	selectProfileQuery := `SELECT id FROM profiles WHERE user_id=?`
	if err := user.DB.Raw(selectProfileQuery, userId).Scan(&profileId).Error; err != nil {
		return "", err
	}
	return profileId, nil
}
func (user *UserAdapter) AdminAddCategory(category entities.Category) error {
	insertCategoryQuery := `INSERT INTO categories (name) VALUES ($1)`
	if err := user.DB.Exec(insertCategoryQuery, category.Name).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) AdminUpdateCategory(category entities.Category) error {
	updateCategory := `UPDATE categories SET name=$1 where id=$2`
	if err := user.DB.Exec(updateCategory, category.Name, category.ID).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) GetAllCategory() ([]entities.Category, error) {
	var res []entities.Category
	selectCategories := `SELECT * FROM categories`
	if err := user.DB.Raw(selectCategories).Scan(&res).Error; err != nil {
		return []entities.Category{}, err
	}
	return res, nil
}
func (user *UserAdapter) GetCategoryByName(category string) (entities.Category, error) {
	var res entities.Category
	selectCategory := `SELECT * FROM categories WHERE name=?`
	if err := user.DB.Raw(selectCategory, category).Scan(&res).Error; err != nil {
		return entities.Category{}, err
	}
	return res, nil
}
func (user *UserAdapter) AdminAddSkill(skill entities.Skill) error {
	var id int
	selectMaxId := `SELECT COALESCE(MAX(id),0) FROM skills`
	if err := user.DB.Raw(selectMaxId).Scan(&id).Error; err != nil {
		return err
	}
	insertSkillQuery := `INSERT INTO skills (id,name,category_id) VALUES ($1,$2,$3)`
	if err := user.DB.Exec(insertSkillQuery, id+1, skill.Name, skill.CategoryId).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) GetSkillByName(skill string) (entities.Skill, error) {
	var res entities.Skill
	selectQuery := `SELECT * FROM skills WHERE name=?`
	if err := user.DB.Raw(selectQuery, skill).Scan(&res).Error; err != nil {
		return entities.Skill{}, err
	}
	return res, nil
}
func (user *UserAdapter) AdminUpdateSkill(skill entities.Skill) error {
	updateskillQuery := `UPDATE skills SET name=$1,category_id=$2 WHERE id=$3`
	if err := user.DB.Exec(updateskillQuery, skill.Name, skill.CategoryId, skill.ID).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) AdminGetAllSkills() ([]helperstruct.SkillHelper, error) {
	var res []helperstruct.SkillHelper
	selectSkillQuery := `SELECT s.id as skill_id,s.name AS skill_name,c.id AS category_id,c.name as category_name FROM skills s JOIN categories c ON c.id=s.category_id`
	if err := user.DB.Raw(selectSkillQuery).Scan(&res).Error; err != nil {
		return []helperstruct.SkillHelper{}, err
	}
	return res, nil

}
func (user *UserAdapter) UserAddSkill(skills entities.UserSkill) error {
	id := uuid.New()
	insertSkillQuery := `INSERT INTO user_skills(id,skill_id,profile_id) VALUES($1,$2,$3)`
	if err := user.DB.Exec(insertSkillQuery, id, skills.SkillId, skills.ProfileId).Error; err != nil {
		return err
	}
	return nil

}
func (user *UserAdapter) UserDeleteSkill(skill entities.UserSkill) error {
	deleteSkillQuery := `DELETE FROM user_skills WHERE skill_id=$1 AND profile_id=$2`
	if err := user.DB.Exec(deleteSkillQuery, skill.SkillId, skill.ProfileId).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) UserGetAllSkills(profileId string) ([]helperstruct.SkillHelper, error) {
	var res []helperstruct.SkillHelper
	selectSkillQueryUser := `SELECT s.id as skill_id,s.name AS skill_name,c.id AS category_id,c.name as category_name FROM skills s JOIN categories c ON c.id=s.category_id JOIN user_skills u ON u.skill_id=s.id WHERE profile_id=$1`
	if err := user.DB.Raw(selectSkillQueryUser, profileId).Scan(&res).Error; err != nil {
		return []helperstruct.SkillHelper{}, err
	}
	return res, nil
}
func (user *UserAdapter) AddLink(links entities.Link) error {
	id := uuid.New()
	insertLinkQuery := `INSERT INTO links (id,profile_id,url,title) VALUES ($1,$2,$3,$4)`
	if err := user.DB.Exec(insertLinkQuery, id, links.ProfileId, links.Url, links.Title).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) GetLinkByTitle(profileId string, title string) (entities.Link, error) {
	selectQuery := `SELECT * FROM links WHERE profile_id=$1 AND title=$2`
	var res entities.Link
	if err := user.DB.Raw(selectQuery, profileId, title).Scan(&res).Error; err != nil {
		return entities.Link{}, err
	}
	return res, nil
}
func (user *UserAdapter) DeleteLink(id string) error {
	deleteQuery := `DELETE FROM links WHERE id=?`
	if err := user.DB.Exec(deleteQuery, id).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) GetAllLinksUser(profileID string) ([]entities.Link, error) {
	var res []entities.Link
	selectLinkQuery := `SELECT * FROM links WHERE profile_id=?`
	if err := user.DB.Raw(selectLinkQuery, profileID).Scan(&res).Error; err != nil {
		return []entities.Link{}, err
	}
	return res, nil
}
func (user *UserAdapter) GetSkillById(id int) (helperstruct.SkillHelper, error) {
	selectSkillQuery := `SELECT s.id AS skill_id,s.name AS skill_name,c.id AS category_id,c.name AS category_name from skills s LEFT JOIN categories c ON c.id=s.category_id where s.id=?`
	var res helperstruct.SkillHelper
	if err := user.DB.Raw(selectSkillQuery, id).Scan(&res).Error; err != nil {
		return helperstruct.SkillHelper{}, err
	}
	return res, nil
}
func (user *UserAdapter) GetCategoryById(id int) (entities.Category, error) {
	selectCategoryQuery := `SELECT * FROM categories WHERE id=?`
	var res entities.Category
	if err := user.DB.Raw(selectCategoryQuery, id).Scan(&res).Error; err != nil {
		return entities.Category{}, err
	}
	return res, nil
}
func (user *UserAdapter) GetUserById(userId string) (entities.User, error) {
	selectUserByIdQuery := `SELECT * from users WHERE id=?`
	var res entities.User
	if err := user.DB.Raw(selectUserByIdQuery, userId).Scan(&res).Error; err != nil {
		return entities.User{}, err
	}
	return res, nil
}
func (user *UserAdapter) JobApply(req entities.Jobs) error {
	id := uuid.New()
	insertIntoQuery := `INSERT INTO jobs (id,job_id,user_id,job_status_id,weightage) VALUES ($1,$2,$3,$4,$5)`
	if err := user.DB.Exec(insertIntoQuery, id, req.JobId, req.UserId, 1, req.Weightage).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) UserEditName(req entities.User) error {
	updateQuery := `UPDATE users SET name=$1 WHERE id=$2`
	if err := user.DB.Exec(updateQuery, req.Name, req.ID).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) UserEditPhone(req entities.User) error {
	updateQuery := `UPDATE users SET phone=$1 WHERE id=$2`
	if err := user.DB.Exec(updateQuery, req.Phone, req.ID).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) UserAddAddress(req entities.Address) error {
	id := uuid.New()
	insertQuery := `INSERT INTO addresses (id,country,state,district,city,profile_id) VALUES ($1,$2,$3,$4,$5,$6)`
	if err := user.DB.Exec(insertQuery, id, req.Country, req.State, req.District, req.City, req.ProfileId).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) UserEditAddress(req entities.Address) error {
	updateQuery := `UPDATE addresses SET country=$1,state=$2,district=$3,city=$4 WHERE profile_id=$5`
	if err := user.DB.Exec(updateQuery, req.Country, req.State, req.District, req.City, req.ProfileId).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) GetAddressByProfileId(id string) (entities.Address, error) {
	var res entities.Address
	selectQuery := `SELECT * FROM addresses WHERE profile_id=?`
	if err := user.DB.Raw(selectQuery, id).Scan(&res).Error; err != nil {
		return entities.Address{}, err
	}
	return res, nil
}
func (user *UserAdapter) UploadProfileImage(image, profileId string) (string, error) {
	var res string
	insertImageQuery := `UPDATE profiles SET image=$1 WHERE id=$2 RETURNING image`
	if err := user.DB.Raw(insertImageQuery, image, profileId).Scan(&res).Error; err != nil {
		return "", err
	}
	return res, nil
}
func (user *UserAdapter) GetProfilePic(profileId string) (string, error) {
	var res string
	selectQuery := `SELECT image FROM profiles WHERE id=$1 AND image NOT IN (NULL)`
	if err := user.DB.Raw(selectQuery, profileId).Scan(&res).Error; err != nil {
		return "", err
	}
	return res, nil
}
func (user *UserAdapter) GetAppliedJobs(userId string) ([]helperstruct.JobHelper, error) {
	var res []helperstruct.JobHelper
	selectQuery := `SELECT j.job_id,js.status FROM jobs j JOIN job_statuses js ON j.job_status_id=js.id WHERE user_id=$1`
	if err := user.DB.Raw(selectQuery, userId).Scan(&res).Error; err != nil {
		return []helperstruct.JobHelper{}, err
	}
	return res, nil
}
func (user *UserAdapter) GetAppliedJob(userId, jobId string) (entities.Jobs, error) {
	var res entities.Jobs
	selectQuery := `SELECT * FROM jobs WHERE job_id=$1 AND user_id=$2`
	if err := user.DB.Raw(selectQuery, jobId, userId).Scan(&res).Error; err != nil {
		return entities.Jobs{}, err
	}
	return res, nil
}
func (user *UserAdapter) GetUserSkillById(profileId string, skillId int) (entities.UserSkill, error) {
	var res entities.UserSkill
	selectQuery := `SELECT * FROM user_skills WHERE profile_id=$1 AND skill_id=$2`
	if err := user.DB.Raw(selectQuery, profileId, skillId).Scan(&res).Error; err != nil {
		return entities.UserSkill{}, err
	}
	return res, nil
}
func (user *UserAdapter) GetAppliedUsersByJobId(jobId string) ([]entities.Jobs, error) {
	selectQuery := `SELECT * FROM jobs WHERE job_id=? ORDER BY weightage DESC`
	var res []entities.Jobs
	if err := user.DB.Raw(selectQuery, jobId).Scan(&res).Error; err != nil {
		return []entities.Jobs{}, err
	}
	return res, nil
}
func (user *UserAdapter) AddExperience(userId, experience string) error {
	updateProfileQuery := `UPDATE profiles SET experience_in_current_field=$1 WHERE user_id=$2`
	if err := user.DB.Exec(updateProfileQuery, experience, userId).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) GetExperience(userId string) (string, error) {
	selectExperience := `SELECT COALESCE(experience_in_current_field,'0 years') FROM profiles WHERE user_id=?`
	var res string
	if err := user.DB.Raw(selectExperience, userId).Scan(&res).Error; err != nil {
		return "", err
	}
	return res, nil
}
func (user *UserAdapter) AddToShortlist(req entities.Shortlist) error {
	insertQuery := `INSERT INTO shortlists (id,user_id,job_id,weightage) VALUES ($1,$2,$3,$4)`
	id := uuid.New()
	if err := user.DB.Exec(insertQuery, id, req.UserId, req.JobId, req.Weightage).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) GetWeightage(userId, jobId string) (float64, error) {
	var res float64
	selectQuery := `SELECT weightage FROM jobs WHERE user_id=$1 AND job_id=$2`
	if err := user.DB.Raw(selectQuery, userId, jobId).Scan(&res).Error; err != nil {
		return 0, err
	}
	return res, nil
}
func (user *UserAdapter) GetShortlist(jobId string) ([]entities.Shortlist, error) {
	var res []entities.Shortlist
	selectQuery := `SELECT * FROM shortlists WHERE job_id=?`
	if err := user.DB.Raw(selectQuery, jobId).Scan(&res).Error; err != nil {
		return []entities.Shortlist{}, err
	}
	return res, nil
}
func (user *UserAdapter) AddEducation(req entities.Education) error {
	id := uuid.New()
	insertQuery := `INSERT INTO educations (id,degree,institution,start_date,end_date,user_id) VALUES ($1,$2,$3,$4,$5,$6)`
	if err := user.DB.Exec(insertQuery, id, req.Degree, req.Institution, req.StartDate, req.EndDate, req.UserId).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) EditEducation(req entities.Education) error {
	updateQuery := `UPDATE educations SET degree=$1,institution=$2,start_date=$3,end_date=$4 WHERE id=$5`
	if err := user.DB.Exec(updateQuery, req.Degree, req.Institution, req.StartDate, req.EndDate, req.ID).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) GetShortlistByUserIdAndJobId(userId, jobId string) (entities.Shortlist, error) {
	var res entities.Shortlist
	selectQuery := `SELECT * FROM shortlists WHERE user_id=$1 AND job_id=$2`
	if err := user.DB.Raw(selectQuery, userId, jobId).Scan(&res).Error; err != nil {
		return entities.Shortlist{}, err
	}
	return res, nil
}
func (user *UserAdapter) GetEducation(userId string) ([]entities.Education, error) {
	selectQuery := `SELECT * FROM educations WHERE user_id=?`
	var res []entities.Education
	if err := user.DB.Raw(selectQuery, userId).Scan(&res).Error; err != nil {
		return []entities.Education{}, err
	}
	return res, nil
}
func (user *UserAdapter) AddToBlockList(userId string) error {
	updateQuery := `UPDATE users SET is_blocked=true WHERE id=?`
	if err := user.DB.Exec(updateQuery, userId).Error; err != nil {
		return err
	}
	return nil

}
func (user *UserAdapter) RemoveFromBlockList(userId string) error {
	updateQuery := `UPDATE users SET is_blocked=false WHERE id=?`
	if err := user.DB.Exec(updateQuery, userId).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) DeleteEducation(edId string) error {
	deleteQuery := `DELETE FROM educations WHERE id=?`
	if err := user.DB.Exec(deleteQuery, edId).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) UpdateAppliedJobStatus(statusId int, jobId, userId string) error {
	updateQuery := `UPDATE jobs SET job_status_id=$1 WHERE user_id=$2 AND job_id=$3`
	if err := user.DB.Exec(updateQuery, statusId, userId, jobId).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) InterviewSchedule(userId, jobId string, date time.Time) error {
	scheduleInterviewQuery := `UPDATE shortlists SET interview_date=$1 WHERE user_id=$2 AND job_id=$3`
	if err := user.DB.Exec(scheduleInterviewQuery, date, userId, jobId).Error; err != nil {
		return err
	}
	return nil
}
func (user *UserAdapter) GetInterviewsForUser(userId string) ([]entities.Shortlist, error) {
	var res []entities.Shortlist
	selectInterviews := `SELECT * FROM shortlists WHERE interview_date IS NOT NULL AND user_id=$1`
	if err := user.DB.Raw(selectInterviews, userId).Scan(&res).Error; err != nil {
		return []entities.Shortlist{}, err
	}
	return res, nil

}
func (user *UserAdapter) ReportUser(userId string) error {
	var reportCount int
	selectReortCount := `SELECT COALESCE(report_count,0) FROM users WHERE id=?`
	if err := user.DB.Raw(selectReortCount, userId).Scan(&reportCount).Error; err != nil {
		return err
	}
	updateQuery := `UPDATE users SET report_count=$1 WHERE id=$2`
	if err := user.DB.Exec(updateQuery, reportCount+1, userId).Error; err != nil {
		return err
	}
	return nil
}
