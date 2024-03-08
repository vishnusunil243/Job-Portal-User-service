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
	ID        uuid.UUID `gorm:"primaryKey;unique;not null"`
	UserId    uuid.UUID
	User      User `gorm:"foreignKey:UserId"`
	Image     string
	Resume    string
	SkillId   uuid.UUID
	Skill     Skill `gorm:"foreignKey:SkillId"`
	LinkId    uuid.UUID
	Link      Link `gorm:"foreignKey:LinkId"`
	AddressId uuid.UUID
	Address   Address `gorm:"foreignKey:AddressId"`
}
type Skill struct {
	ID   uuid.UUID `gorm:"primaryKey;unique;not null"`
	Name string
}
type Link struct {
	ID    uuid.UUID `gorm:"primaryKey;unique;not null"`
	Title uuid.UUID
	Url   string
}
type Address struct {
	Id       uuid.UUID `gorm:"primaryKey;unique;not null"`
	Country  string
	State    string
	District string
	City     string
}
