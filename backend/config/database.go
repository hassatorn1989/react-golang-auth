package config

import (
	"auth-backend/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(cfg AppConfig) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.User{}, &models.RefreshToken{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
