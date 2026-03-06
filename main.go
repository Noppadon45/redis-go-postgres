package main

import (
	"log"

	"github.com/Noppadon/db"
	"github.com/Noppadon/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Initialize Postgres Database
	db.ConnectDB()

	// Setup Routes
	app.Post("/users", handlers.CreateUser)
	app.Get("/users/:id", handlers.GetUser)
	app.Put("/users/:id", handlers.UpdateUser)
	app.Delete("/users/:id", handlers.DeleteUser)

	log.Fatal(app.Listen(":3000"))
}
