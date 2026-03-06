package db

import (
	"fmt"
	"log"

	"github.com/Noppadon/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=admin password=admin1234 dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	fmt.Println("Connected to Database")

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database. \n", err)
	}
	fmt.Println("Database Migrated")

	DB = db
}
