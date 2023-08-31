package controllers

import (
	"avito-task/entity"
	"avito-task/services"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddUserHandler(c *fiber.Ctx) error {
	var payload entity.UserCreateSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	errors := entity.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	user, err := services.Users.AddOne(payload.ID)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User with such ID already exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": fiber.Map{"user": user}})
}

func GetUsersHandler(c *fiber.Ctx) error {
	users, err := services.Users.GetAll()
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	var filteredUsers = make([]*entity.UserWithSegmentsSchema, len(users))
	for i, user := range users {
		f := &entity.UserWithSegmentsSchema{ID: user.ID}
		for _, segment := range user.Segments {
			f.Segments = append(f.Segments, segment.Name)
		}
		filteredUsers[i] = f
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"results": len(users), "data": filteredUsers})
}

func GetUserHandler(c *fiber.Ctx) error {
	var payload entity.UserCreateSchema
	if err := c.ParamsParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	errors := entity.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	user, err := services.Users.GetOne(payload.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User with such ID was not found"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"user": user}})
}
