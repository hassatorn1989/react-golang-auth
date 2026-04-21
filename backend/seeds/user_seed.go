package seeds

import (
	"auth-backend/utils"
	"fmt"
	"log"

	"auth-backend/models"

	"gorm.io/gorm"
)

type SeedUser struct {
	Name     string
	Email    string
	Password string
}

func SeedUsers(db *gorm.DB) {
	users := []SeedUser{
		{Name: "Admin", Email: "admin@example.com", Password: "123456"},
		{Name: "User", Email: "user@example.com", Password: "123456"},
	}

	for _, item := range users {
		var existing models.User
		err := db.Where("email = ?", item.Email).First(&existing).Error
		if err == nil {
			fmt.Printf("skip: %s already exists\n", item.Email)
			continue
		}

		hashedPassword, err := utils.HashPassword(item.Password)
		if err != nil {
			log.Fatalf("hash error for %s: %v", item.Email, err)
		}

		user := models.User{
			Name:         item.Name,
			Email:        item.Email,
			PasswordHash: hashedPassword,
		}

		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("create error for %s: %v", item.Email, err)
		}

		fmt.Printf("created: %s\n", item.Email)
	}
}
