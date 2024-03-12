package entities

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;unique;not null"`
	Name     string
	Email    string
	Phone    string
	Password string
}
type Profile struct {
	ID     uuid.UUID `gorm:"primaryKey;unique;not null"`
	UserId uuid.UUID
	User   User `gorm:"foreignKey:UserId"`
	Image  string
	Resume string
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
