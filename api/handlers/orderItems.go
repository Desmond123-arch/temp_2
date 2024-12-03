package handlers

import (
	"log"

	"devops.zedeks.com/TheHiddenDeveloper/ims-zedeks/api/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllOrderItems(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}
	
	_, err := db.NewCreateTable().Model(&models.OrderItem{}).IfNotExists().Exec(dbCtx)
	
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "An unexpected error occurred",
		})
	}

	var orderItems []models.OrderItem
	err = db.NewSelect().Model(&orderItems).Scan(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	if len(orderItems) == 0 {
		return c.Status(fiber.StatusNoContent).JSON([]models.OrderItem{})
	}
	return c.Status(fiber.StatusOK).JSON(orderItems)
}

func CreateOrderItem(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	var orderItem models.OrderItem
	if err := c.BodyParser(&orderItem); err != nil {
		log.Printf("Parse Error: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	_, err = db.NewInsert().Model(&orderItem).Exec(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create order item",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(orderItem)
}

func GetOneOrderItem(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	id := c.Params("id")
	var orderItem models.OrderItem

	err = db.NewSelect().Model(&orderItem).Where("id = ?", id).Scan(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order item not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(orderItem)
}

func UpdateOrderItem(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	id := c.Params("id")
	var orderItem models.OrderItem

	err = db.NewSelect().Model(&orderItem).Where("id = ?", id).Scan(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order item not found",
		})
	}

	if err := c.BodyParser(&orderItem); err != nil {
		log.Printf("Parse Error: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	_, err = db.NewUpdate().Model(&orderItem).Where("id = ?", id).Exec(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order item",
		})
	}

	return c.Status(fiber.StatusOK).JSON(orderItem)
}

func DeleteOrderItem(c *fiber.Ctx) error {
	if err != nil {
		log.Printf("Database Error: %s", err)
		return err
	}

	id := c.Params("id")

	result, err := db.NewDelete().Model((*models.OrderItem)(nil)).Where("id = ?", id).Exec(dbCtx)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete order item",
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order item not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
