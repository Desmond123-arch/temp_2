package handlers

import (
	"fmt"
	"log"
	"strings"

	"devops.zedeks.com/TheHiddenDeveloper/ims-zedeks/api/models"
	"github.com/gofiber/fiber/v2"
)

// Get all suppliers
func GetAllSuppliers(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}
	
	_, err := db.NewCreateTable().Model(&models.Supplier{}).IfNotExists().Exec(dbCtx)
	
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "An unexpected error occurred",
		})
	}

	var suppliers []models.Supplier
	err = db.NewSelect().Model(&suppliers).Scan(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	if len(suppliers) == 0 {
		return c.Status(fiber.StatusNoContent).JSON([]models.Supplier{})
	}
	return c.Status(fiber.StatusOK).JSON(suppliers)
}

func CreateSupplier(c *fiber.Ctx) error {
	var supplier models.Supplier
	if err := c.BodyParser(&supplier); err != nil {
		log.Printf("Parse Error: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if supplier.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": "Supplier name is required and cannot be empty",
		})
	}

	// Check if supplier with same name or email already exists
	var existingSupplier models.Supplier
	err := db.NewSelect().
		Model(&existingSupplier).
		Where("LOWER(name) = LOWER(?)", supplier.Name).
		WhereOr("LOWER(email) = LOWER(?)", supplier.Email).
		Scan(dbCtx)
	
	if err == nil {
		// Determine which field caused the conflict
		field := "name"
		if strings.EqualFold(existingSupplier.Email, supplier.Email) {
			field = "email"
		}
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Duplicate entry",
			"details": fmt.Sprintf("A supplier with this %s already exists", field),
			"field": field,
		})
	}

	_, err = db.NewInsert().Model(&supplier).Exec(dbCtx)
	if err != nil {
		log.Printf("Full Database Error: %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database operation failed",
			"details": "Failed to create supplier. Please try again later",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(supplier)
}

func GetOneSupplier(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	id := c.Params("id")
	var supplier models.Supplier

	err = db.NewSelect().Model(&supplier).Where("id = ?", id).Scan(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Supplier not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(supplier)
}

func UpdateSupplier(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	id := c.Params("id")
	var originalSupplier models.Supplier

	// Check if supplier exists
	err = db.NewSelect().Model(&originalSupplier).Where("id = ?", id).Scan(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Supplier not found",
		})
	}

	var supplier models.Supplier
	if err := c.BodyParser(&supplier); err != nil {
		log.Printf("Parse Error: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate required fields
	if supplier.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": "Supplier name is required and cannot be empty",
		})
	}

	_, err = db.NewUpdate().Model(&supplier).Where("id = ?", id).Exec(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update supplier",
		})
	}

	return c.Status(fiber.StatusOK).JSON(supplier)
}

func DeleteSupplier(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	id := c.Params("id")

	result, err := db.NewDelete().Model((*models.Supplier)(nil)).Where("id = ?", id).Exec(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete supplier",
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Supplier not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}