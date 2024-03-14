package db

import (
	"github.com/vishnusunil243/Job-Portal-User-service/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(connectTo string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectTo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&entities.User{})
	db.AutoMigrate(&entities.Profile{})
	db.AutoMigrate(&entities.Address{})
	db.AutoMigrate(&entities.Link{})
	db.AutoMigrate(&entities.Skill{})
	db.AutoMigrate(&entities.Admin{})
	db.AutoMigrate(&entities.UserSkill{})
	db.AutoMigrate(&entities.Jobs{})
	db.AutoMigrate(&entities.JobStatus{})
	return db, nil
}
