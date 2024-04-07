package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"primaryKey;unique;not null"`
	Name        string
	Email       string
	Phone       string
	Password    string
	IsBlocked   bool
	ReportCount int
	CreatedAt   time.Time
	Subscribed  bool
}
type Profile struct {
	ID                       uuid.UUID `gorm:"primaryKey;unique;not null"`
	UserId                   uuid.UUID
	User                     User `gorm:"foreignKey:UserId"`
	Image                    string
	Resume                   string
	ExperienceInCurrentField string
}
type Category struct {
	ID   int
	Name string
}
type Skill struct {
	ID         int `gorm:"primaryKey;unique;not null"`
	CategoryId int
	Category   Category `gorm:"foreignKey:CategoryId"`
	Name       string
}
type UserSkill struct {
	ID        uuid.UUID
	ProfileId uuid.UUID
	Profile   Profile `gorm:"foreignKey:ProfileId"`
	SkillId   int
	Skill     Skill `gorm:"foreignKey:SkillId"`
}
type Link struct {
	ID        uuid.UUID `gorm:"primaryKey;unique;not null"`
	Title     string
	ProfileId uuid.UUID
	Profile   Profile `gorm:"foreignKey:ProfileId"`
	Url       string
}
type Address struct {
	Id        uuid.UUID `gorm:"primaryKey;unique;not null"`
	Country   string
	State     string
	District  string
	City      string
	ProfileId uuid.UUID
	Profile   Profile `gorm:"foreignKey:ProfileId"`
}
type Admin struct {
	ID       uuid.UUID
	Name     string
	Password string
	Email    string
	Phone    string
}
type Jobs struct {
	Id            uuid.UUID
	UserId        uuid.UUID
	JobId         string
	User          User `gorm:"foreignKey:UserId"`
	JobStatusId   int
	JobStatus     JobStatus `gorm:"foreignKey:JobStatusId"`
	Weightage     float64
	InterviewDate time.Time
}
type JobStatus struct {
	Id     int
	Status string
}
type Shortlist struct {
	ID            uuid.UUID
	UserId        uuid.UUID
	JobId         uuid.UUID
	Weightage     float64
	InterviewDate time.Time
	RoomId        string
	Status        string
	WarningSent   bool
}

type Education struct {
	ID          uuid.UUID
	UserId      uuid.UUID
	User        User `gorm:"foreignKey:UserId"`
	Degree      string
	Institution string
	StartDate   time.Time `gorm:"type:date"`
	EndDate     time.Time `gorm:"type:date"`
}
type Project struct {
	Id          uuid.UUID
	UserId      uuid.UUID
	User        User `gorm:"foreignKey:UserId"`
	Name        string
	Description string
	Image       string
	Link        string
}
