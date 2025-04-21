package database

import (
	"server/internal/models"

	// "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// database, err := gorm.Open(mysql.Open("go_admin:password@/go_admin"), &gorm.Config{})
	database, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		panic("Could not connect to the database")
	}

	DB = database

	database.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})
}
