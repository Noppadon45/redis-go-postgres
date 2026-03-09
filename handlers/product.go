package handlers

import (
	"context"
	"encoding/json"

	"github.com/Noppadon/db"
	"github.com/Noppadon/models"
	"github.com/gofiber/fiber/v3"
)

// CreateProduct handles POST /products
func CreateProduct(c fiber.Ctx) error {
	var product models.Product
	if err := c.Bind().Body(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	result := db.DB.Create(&product)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}

	// Invalidate Cache
	db.RedisDB.HDel(context.Background(), "product", "data")

	return c.Status(201).JSON(product)
}

// GetAllProducts handles GET /products
func GetAllProducts(c fiber.Ctx) error {
	var products []models.Product

	// 1. Check Redis HGet
	val, err := db.RedisDB.HGet(context.Background(), "product", "data").Result()
	if err == nil {
		// Found in cache, return it
		if jsonErr := json.Unmarshal([]byte(val), &products); jsonErr == nil {
			return c.JSON(products)
		}
	}

	// Preload "User" to attach user details to the product response
	result := db.DB.Preload("User").Find(&products)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}

	// 2. Set HSet in Redis
	if productsJSON, err := json.Marshal(products); err == nil {
		db.RedisDB.HSet(context.Background(), "product", "data", productsJSON)
	}

	return c.JSON(products)
}

// GetProduct handles GET /products/:id
func GetProduct(c fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	// Preload "User" to attach user details
	result := db.DB.Preload("User").First(&product, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}
	return c.JSON(product)
}

// UpdateProduct handles PUT /products/:id
func UpdateProduct(c fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	if err := db.DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	type UpdateProductInput struct {
		Title  string `json:"title"`
		UserID uint   `json:"user_id"`
	}
	var input UpdateProductInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	db.DB.Model(&product).Updates(models.Product{
		Title:  input.Title,
		UserID: input.UserID,
	})

	// Fetch updated product with user details
	db.DB.Preload("User").First(&product, id)

	// Invalidate Cache
	db.RedisDB.HDel(context.Background(), "product", "data")

	return c.JSON(product)
}

// DeleteProduct handles DELETE /products/:id
func DeleteProduct(c fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	if err := db.DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	db.DB.Delete(&product)

	// Invalidate Cache
	db.RedisDB.HDel(context.Background(), "product", "data")

	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}
