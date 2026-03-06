package handlers

import (
	"github.com/Noppadon/db"
	"github.com/Noppadon/models"
	"github.com/gofiber/fiber/v2"
)

// CreateUser handles POST /users
func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.Status(201).JSON(user)
}

// GetUser handles GET /users/:id
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	result := db.DB.First(&user, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

// UpdateUser handles PUT /users/:id
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := db.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	type UpdateUserInput struct {
		Fullname string `json:"Fullname"`
		Age      string `json:"Age"`
		Location string `json:"Location"`
		Email    string `json:"Email"`
		Zipcode  string `json:"Zipcode"`
	}
	var input UpdateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	db.DB.Model(&user).Updates(models.User{
		Fullname: input.Fullname,
		Age:      input.Age,
		Location: input.Location,
		Email:    input.Email,
		Zipcode:  input.Zipcode,
	})

	return c.JSON(user)
}

// DeleteUser handles DELETE /users/:id
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := db.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	db.DB.Delete(&user)
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}
