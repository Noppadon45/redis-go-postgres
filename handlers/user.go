package handlers

import (
	"github.com/Noppadon/db"
	"github.com/Noppadon/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser handles POST /users
func CreateUser(c fiber.Ctx) error {
	type CreateUserInput struct {
		Fullname string `json:"Fullname"`
		Age      string `json:"Age"`
		Location string `json:"Location"`
		Email    string `json:"Email"`
		Password string `json:"Password"`
		Zipcode  string `json:"Zipcode"`
	}
	var input CreateUserInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
	}

	user := models.User{
		Fullname: input.Fullname,
		Age:      input.Age,
		Location: input.Location,
		Email:    input.Email,
		Password: string(hashedPassword),
		Zipcode:  input.Zipcode,
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.Status(201).JSON(user)
}

// GetAllUsers handles GET /users
func GetAllUsers(c fiber.Ctx) error {
	var users []models.User
	result := db.DB.Find(&users)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.JSON(users)
}

// GetUser handles GET /users/:id
func GetUser(c fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	result := db.DB.First(&user, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

// UpdateUser handles PUT /users/:id
func UpdateUser(c fiber.Ctx) error {
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
		Password string `json:"Password"`
		Zipcode  string `json:"Zipcode"`
	}
	var input UpdateUserInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	updateData := models.User{
		Fullname: input.Fullname,
		Age:      input.Age,
		Location: input.Location,
		Email:    input.Email,
		Zipcode:  input.Zipcode,
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err == nil {
			updateData.Password = string(hashedPassword)
		}
	}

	db.DB.Model(&user).Updates(updateData)

	return c.JSON(user)
}

// DeleteUser handles DELETE /users/:id
func DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := db.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	db.DB.Delete(&user)
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

func Login(c fiber.Ctx) error {
	// 1. รับ Email และ Password จากผู้ใช้
	var input struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}
	// ... parse input ...
	// 2. ดึง User จาก Database ด้วย Email
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	// 3.เปรียบเทียบ Password ดิบ กับ Hash ในระบบ 
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		// ถ้ารหัสไม่ตรงกัน มันจะ return error คืนมา
		return c.Status(401).JSON(fiber.Map{"error": "Invalid password"})
	}

	// 4.สร้าง Session และเก็บข้อมูลผู้ใช้ 
	sess := session.FromContext(c)
	if sess == nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to get session from context"})
	}

	// Important: Regenerate the session ID to prevent fixation
	// This changes the session ID while preserving existing data
	if err := sess.Regenerate(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to regenerate session"})
	}

	sess.Set("Fullname", user.Fullname)
	sess.Set("Age", user.Age)
	sess.Set("Location", user.Location)
	sess.Set("Zipcode", user.Zipcode)

	// ผ่าน! -> ออก Token หรือ อนุญาตให้เข้าระบบ
	return c.Status(200).JSON(fiber.Map{
		"message":  "Login successful",
		"Fullname": user.Fullname,
		"Age":      user.Age,
		"Location": user.Location,
		"Zipcode":  user.Zipcode,
	})
}

// Logout handles POST /logout
func Logout(c fiber.Ctx) error {
	sess := session.FromContext(c)
	if sess == nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to get session from context"})
	}

	// Complete session reset (clears all data + new session ID)
	if err := sess.Reset(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to reset session"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Logout successful"})
}
