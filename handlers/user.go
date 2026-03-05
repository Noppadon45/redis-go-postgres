package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type User struct {
	Fullname string `redis:"Fullname"`
	Age      string `redis:"Age"`
	Location string `redis:"Location"`
	Email    string `redis:"Email"`
	Zipcode  string `redis:"Zipcode"`
}

func CreateUser(rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()

		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		err := rdb.HSet(ctx, "user:"+ user.Email, user).Err()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "User created"})

	}
}

func GetUser(rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		id := c.Params("id")

		var user User
		err := rdb.HGetAll(ctx, "user:"+id).Scan(&user)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(user)
	}
}
