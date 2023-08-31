package controllers

import (
	"avito-task/entity"
	"avito-task/initializers"
	"avito-task/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetSegmentLogsHandler(c *fiber.Ctx) error {
	var payload entity.SegmentLogRequestSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	errors := entity.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	path, err := services.SegmentLogs.GenerateCSV(payload.UserID, payload.DateAfter, payload.DateBefore)
	fmt.Println("Works 1")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println("Works2")
	config, _ := initializers.LoadConfig(".")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": config.StaticPath + path,
	})
}
