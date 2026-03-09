package db

import (
	"fmt"
	"log"

	"github.com/Noppadon/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisDB *redis.Client

func ConnectDB() {
	dsn := "host=localhost user=admin password=admin1234 dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	RedisDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	fmt.Println("Connected to Database and Redis")

	err = db.AutoMigrate(&models.User{}, &models.Product{})
	if err != nil {
		log.Fatal("Failed to migrate database. \n", err)
	}
	fmt.Println("Database Migrated")

	DB = db
}
