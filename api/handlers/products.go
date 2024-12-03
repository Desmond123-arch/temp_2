package handlers

import (
	"context"
	"log"
	"time"

	"devops.zedeks.com/TheHiddenDeveloper/ims-zedeks/api/database"
	"devops.zedeks.com/TheHiddenDeveloper/ims-zedeks/api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var db, err = database.ConnectDb()
var dbCtx context.Context
var cancel context.CancelFunc
func init() {
	dbCtx, cancel = context.WithTimeout(context.Background(), 30*24*time.Hour)
}

func Getall(c *fiber.Ctx) error {
	log.Println("Starting Getall function")
	
	if err != nil {
		log.Printf("Initial Database Connection Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection failed",
		})
	}
	log.Println("Database connection successful")
	
	// Create table if not exists
	log.Println("Attempting to create table if not exists")
	_, err := db.NewCreateTable().
		Model(&models.Products{}).
		IfNotExists().
		Exec(dbCtx)
	
	if err != nil {
		log.Printf("Table Creation Error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create table",
			"details": err.Error(),
		})
	}
	log.Println("Table creation/check completed")

	log.Println("Attempting to fetch products")
	var products []models.Products
	err = db.NewSelect().
		Model(&products).
		Relation("Category").
		Relation("Supplier").
		
		Scan(dbCtx)

	if err != nil {
		log.Printf("Product Query Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
			"details": err.Error(),
		})
	}
	log.Printf("Successfully fetched %d products", len(products))

	if len(products) == 0 {
		return c.Status(fiber.StatusNoContent).JSON([]models.Products{})
	}
	return c.Status(fiber.StatusOK).JSON(products)
}

// Create a new product
func Create(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	// Create a struct to parse the JSON request
	var requestData struct {
		Name       string  `json:"name"`
		CategoryID string  `json:"category_id"`
		Price      float64 `json:"price"`
		Quantity   int     `json:"quantity"`
		ImageURL   string  `json:"image_url,omitempty"`
		SupplierID string  `json:"supplier_id"`
	}

	// Parse JSON body
	if err := c.BodyParser(&requestData); err != nil {
		log.Printf("Parse Error: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Parse Category ID
	categoryID, err := uuid.Parse(requestData.CategoryID)
	if err != nil {
		log.Printf("Category ID Parse Error: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID format",
			"details": err.Error(),
		})
	}

	// Parse Supplier ID
	supplierID, err := uuid.Parse(requestData.SupplierID)
	if err != nil {
		log.Printf("Supplier ID Parse Error: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid supplier ID format",
			"details": err.Error(),
		})
	}

	// Create the product
	product := models.Products{
		Name:       requestData.Name,
		CategoryID: categoryID,
		Price:      requestData.Price,
		Quantity:   requestData.Quantity,
		ImageURL:   requestData.ImageURL,
		SupplierID: supplierID,
	}

	// Insert the product
	_, err = db.NewInsert().
		Model(&product).
		Exec(dbCtx)

	if err != nil {
		log.Printf("Insert Error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
			"details": err.Error(),
		})
	}

	// Load relations
	err = db.NewSelect().
		Model(&product).
		Relation("Category").
		Relation("Supplier").
		Where("products.id = ?", product.ID).
		Scan(dbCtx)

	if err != nil {
		log.Printf("Error loading relations: %s", err)
	}

	// Create response struct without CategoryID and SupplierID
	response := struct {
		ID       uuid.UUID        `json:"ID"`
		Name     string          `json:"Name"`
		Category models.Category `json:"Category"`
		Price    float64         `json:"Price"`
		Quantity int             `json:"Quantity"`
		ImageURL string          `json:"ImageURL"`
		Supplier models.Supplier `json:"Supplier"`
	}{
		ID:       product.ID,
		Name:     product.Name,
		Category: product.Category,
		Price:    product.Price,
		Quantity: product.Quantity,
		ImageURL: product.ImageURL,
		Supplier: product.Supplier,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// Get a single product by ID
func GetOne(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	id := c.Params("id")
	var product models.Products

	err = db.NewSelect().
		Model(&product).
		Relation("Category").
		Relation("Supplier").
		Where("products.id = ?", id).
		Scan(dbCtx)

	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Create response struct without CategoryID and SupplierID
	response := struct {
		ID       uuid.UUID        `json:"ID"`
		Name     string          `json:"Name"`
		Category models.Category `json:"Category"`
		Price    float64         `json:"Price"`
		Quantity int             `json:"Quantity"`
		ImageURL string          `json:"ImageURL"`
		Supplier models.Supplier `json:"Supplier"`
	}{
		ID:       product.ID,
		Name:     product.Name,
		Category: product.Category,
		Price:    product.Price,
		Quantity: product.Quantity,
		ImageURL: product.ImageURL,
		Supplier: product.Supplier,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Update a product
func Update(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	id := c.Params("id")
	var product models.Products
	err = db.NewSelect().Model(&product).Where("id = ?", id).Scan(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	if err := c.BodyParser(&product); err != nil {
		log.Printf("Parse Error: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	_, err = db.NewUpdate().Model(&product).Where("id = ?", id).Exec(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

// Delete a product
func Delete(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	id := c.Params("id")

	result, err := db.NewDelete().Model((*models.Products)(nil)).Where("id = ?", id).Exec(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product",
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetProductsByCategory retrieves all products in a specific category
func GetProductsByCategory(c *fiber.Ctx) error {
	log.Println("Starting GetProductsByCategory function")
	
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	categoryID := c.Params("categoryId")
	
	// Parse the category ID string to UUID
	parsedCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		log.Printf("Category ID Parse Error: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID format",
			"details": err.Error(),
		})
	}

	// Modified to select only ID and Name
	var products []struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
	err = db.NewSelect().
		Model((*models.Products)(nil)).
		Column("id", "name").
		Where("category_id = ?", parsedCategoryID).
		Scan(dbCtx, &products)

	if err != nil {
		log.Printf("Database Query Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
			"details": err.Error(),
		})
	}

	log.Printf("Successfully fetched %d products for category %s", len(products), categoryID)

	if len(products) == 0 {
		return c.Status(fiber.StatusNoContent).JSON([]struct{}{})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

// GetProductsBySupplier retrieves all products from a specific supplier
func GetProductsBySupplier(c *fiber.Ctx) error {
	log.Println("Starting GetProductsBySupplier function")
	
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	supplierID := c.Params("supplierId")
	
	// Parse the supplier ID string to UUID
	parsedSupplierID, err := uuid.Parse(supplierID)
	if err != nil {
		log.Printf("Supplier ID Parse Error: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid supplier ID format",
			"details": err.Error(),
		})
	}

	// Modified to select only ID and Name
	var products []struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
	err = db.NewSelect().
		Model((*models.Products)(nil)).
		Column("id", "name").
		Where("supplier_id = ?", parsedSupplierID).
		Scan(dbCtx, &products)

	if err != nil {
		log.Printf("Database Query Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
			"details": err.Error(),
		})
	}

	log.Printf("Successfully fetched %d products for supplier %s", len(products), supplierID)

	if len(products) == 0 {
		return c.Status(fiber.StatusNoContent).JSON([]struct{}{})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

