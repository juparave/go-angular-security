package database

import (
	"log"
	"server/internal/models"

	// "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var appModels = []interface{}{
	&models.User{},
	&models.Role{},
	&models.Permission{},
}

var DB *gorm.DB

func Connect() {
	// database, err := gorm.Open(mysql.Open("go_admin:password@/go_admin"), &gorm.Config{})
	database, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		log.Println("Could not connect to the database:", err)
		panic("Could not connect to the database")
	}

	DB = database

	if err := database.AutoMigrate(appModels...); err != nil {
		log.Println("Could not migrate the database:", err)
		panic("Could not migrate the database")
	}
}
