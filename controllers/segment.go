package controllers

import (
	"avito-task/entity"
	"avito-task/services"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddSegmentHandler(c *fiber.Ctx) error {
	var payload entity.SegmentCreateSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	errors := entity.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	var segment *entity.Segment
	var err error
	if payload.Percent != 0 {
		segment, err = services.Segments.AddOneWithPercent(payload.Name, payload.Percent)
	} else {
		segment, err = services.Segments.AddOne(payload.Name)
	}
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Segment with such name already exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": fiber.Map{"segment": segment}})
}

func GetAllSegmentsHandler(c *fiber.Ctx) error {
	segments, err := services.Segments.GetAll()
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	var filteredSegments = make([]*entity.SegmentWithUsersSchema, len(segments))
	for i, s := range segments {
		f := &entity.SegmentWithUsersSchema{Name: s.Name}
		for _, u := range s.Users {
			f.Users = append(f.Users, u.ID)
		}
		filteredSegments[i] = f
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"results": len(segments), "segments": filteredSegments})
}

func GetSegmentHandler(c *fiber.Ctx) error {
	name := c.Params("segment")
	if len(name) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Segment name not specified"})
	}
	segment, err := services.Segments.GetOne(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Segment with such name was not found"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	filteredSegment := &entity.SegmentWithUsersSchema{Name: segment.Name}
	for _, u := range segment.Users {
		filteredSegment.Users = append(filteredSegment.Users, u.ID)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"segment": filteredSegment}})
}

func UpdateSegmentHandler(c *fiber.Ctx) error {
	var payload entity.SegmentUpdateSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	errors := entity.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	segment, err := services.Segments.UpdateOne(payload.Name, payload.NewName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Segment with such name was not found"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": fiber.Map{"segment": segment}})
}

func DeleteSegmentHandler(c *fiber.Ctx) error {
	var payload entity.SegmentCreateSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	errors := entity.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	err := services.Segments.DeleteOne(payload.Name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Segment with such name was not found"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "The segment was deleted successfully"})
}

func AddSegmentsToUserHandler(c *fiber.Ctx) error {
	var payload entity.AddSegmentsSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	names := strings.Split(payload.SegmentNames, ",")
	errors := entity.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	err := services.Segments.AddSegmentsToUser(payload.ID, names)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Successfully added segments to user " + fmt.Sprint(payload.ID)})
}

func RemoveUserFromSegment(c *fiber.Ctx) error {
	var payload entity.RemoveUserFromSegmentSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	errors := entity.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	err := services.Segments.RemoveUserFromSegment(payload.ID, payload.Segment)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Segment or user was not found"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "The segment was deleted successfully"})
}

func AddAndRemoveSegmentsHandler(c *fiber.Ctx) error {
	var payload entity.AddAndRemoveSegmentsSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	errors := entity.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	err := services.Segments.AddSegmentsToUser(payload.ID, payload.AddSegments)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Segment or user was not found"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	err = services.Segments.RemoveSegmentsFromUser(payload.ID, payload.RemoveSegments)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Segment or user was not found"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Successfully added and removed the specified segments"})
}
