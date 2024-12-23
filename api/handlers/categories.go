package handlers

import (
	"fmt"
	"log"
	"strings"

	"devops.zedeks.com/TheHiddenDeveloper/ims-zedeks/api/models"
	"github.com/gofiber/fiber/v2"
)

// Get all categories
func GetAllCategories(c *fiber.Ctx) error {
	_, err := db.NewCreateTable().
		Model(&models.Category{}).
		IfNotExists().
		WithForeignKeys().
		Exec(dbCtx)
	
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "An unexpected error occurred",
		})
	}

	var categories []models.Category
	err = db.NewSelect().Model(&categories).Scan(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	if len(categories) == 0 {
		return c.Status(fiber.StatusNoContent).JSON([]models.Category{})
	}
	return c.Status(fiber.StatusOK).JSON(categories)
}

// Create a new category
func CreateCategory(c *fiber.Ctx) error {
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		log.Printf("Parse Error: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if category.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": "Category name is required and cannot be empty",
		})
	}

	log.Printf("Attempting to create category: %+v", category)

	_, err := db.NewInsert().Model(&category).Exec(dbCtx)
	if err != nil {
		log.Printf("Full Database Error: %+v", err)
		
		// Check for unique constraint violations
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint") || 
		   strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Duplicate entry",
				"details": fmt.Sprintf("A category with the name '%s' already exists", category.Name),
				"field": "name",
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database operation failed",
			"details": "Failed to create category. Please try again later",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// Get a single category by ID
func GetOneCategory(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	id := c.Params("id")
	var category models.Category

	err = db.NewSelect().Model(&category).Where("id = ?", id).Scan(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}

// Update a category
func UpdateCategory(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	id := c.Params("id")
	var originalCategory models.Category

	// Check if category exists
	err = db.NewSelect().Model(&originalCategory).Where("id = ?", id).Scan(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		log.Printf("Parse Error: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if category.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": "Category name is required",
		})
	}

	// Preserve the ID from the original category
	category.ID = originalCategory.ID

	// Perform the update
	_, err = db.NewUpdate().Model(&category).Where("id = ?", id).Exec(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		// Check for unique constraint violations
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint") || 
		   strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Duplicate entry",
				"details": fmt.Sprintf("A category with the name '%s' already exists", category.Name),
				"field": "name",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update category",
			"details": "Database operation failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}

// Delete a category
func DeleteCategory(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	id := c.Params("id")

	result, err := db.NewDelete().Model((*models.Category)(nil)).Where("id = ?", id).Exec(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete category",
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
