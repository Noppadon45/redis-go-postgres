package main

import (
	"log"
	"time"

	"github.com/Noppadon/db"
	"github.com/Noppadon/handlers"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/redis/v3"
)

func main() {
	app := fiber.New()

	// Initialize Postgres Database
	db.ConnectDB()

	// Initialize Redis Storage for Session
	storage := redis.New(redis.Config{
		Host: "localhost",
		Port: 6379,
	})

	app.Use(session.New(session.Config{
		Storage:         storage,
		CookieSecure:    true,             // HTTPS only
		CookieHTTPOnly:  true,             // Prevent XSS
		CookieSameSite:  "Lax",            // CSRF protection
		IdleTimeout:     30 * time.Minute, // Session timeout
		AbsoluteTimeout: 24 * time.Hour,   // Maximum session life
	}))

	// Setup Routes
	app.Post("/users", handlers.CreateUser)
	app.Get("/users", handlers.GetAllUsers)
	app.Get("/users/:id", handlers.GetUser)
	app.Put("/users/:id", handlers.UpdateUser)
	app.Delete("/users/:id", handlers.DeleteUser)

	// Product Routes
	app.Post("/products", handlers.CreateProduct)
	app.Get("/products", handlers.GetAllProducts)
	app.Get("/products/:id", handlers.GetProduct)
	app.Put("/products/:id", handlers.UpdateProduct)
	app.Delete("/products/:id", handlers.DeleteProduct)

	app.Post("/login", handlers.Login)
	app.Post("/logout", handlers.Logout)

	log.Fatal(app.Listen(":3000"))
}
