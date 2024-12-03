package main

import (
	"fmt"
	"os"
	"devops.zedeks.com/TheHiddenDeveloper/ims-zedeks/api/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	app := fiber.New()

	products_endpoints := app.Group("/products")
	
	products_endpoints.Get("/", handlers.Getall)
	products_endpoints.Post("/", handlers.Create)
	products_endpoints.Get("/:id", handlers.GetOne)
	products_endpoints.Put("/:id", handlers.Update)
	products_endpoints.Delete("/:id", handlers.Delete)
	app.Get("/categories/:categoryId/products", handlers.GetProductsByCategory)
	app.Get("/suppliers/:supplierId/products", handlers.GetProductsBySupplier)

	categories_endpoints := app.Group("/categories")
	categories_endpoints.Get("/", handlers.GetAllCategories)
	categories_endpoints.Post("/", handlers.CreateCategory)
	categories_endpoints.Get("/:id", handlers.GetOneCategory)
	categories_endpoints.Put("/:id", handlers.UpdateCategory)
	categories_endpoints.Delete("/:id", handlers.DeleteCategory)

	suppliers_endpoints := app.Group("/suppliers")
	suppliers_endpoints.Get("/", handlers.GetAllSuppliers)
	suppliers_endpoints.Post("/", handlers.CreateSupplier)
	suppliers_endpoints.Get("/:id", handlers.GetOneSupplier)
	suppliers_endpoints.Put("/:id", handlers.UpdateSupplier)
	suppliers_endpoints.Delete("/:id", handlers.DeleteSupplier)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // fallback to 3000 if PORT is not set
	}

	app.Listen(":" + port)
}