package main

import (
	"avito-task/controllers"
	"avito-task/initializers"
	"avito-task/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables\n", err.Error())
	}
	initializers.ConnectDB(&config)
	services.Segments = services.CreateSegmentService(initializers.DB)
	services.Users = services.CreateUserService(initializers.DB)
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/api/healthcheck", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Everything is fine.",
		})
	})
	users := app.Group("/api/users")
	users.Post("/", controllers.AddUserHandler)
	users.Get("/", controllers.GetUsersHandler)
	users.Get("/:id", controllers.GetUserHandler)

	segments := app.Group("/api/segments")
	segments.Get("/", controllers.GetAllSegmentsHandler)
	segments.Get("/:segment", controllers.GetSegmentHandler)
	segments.Post("/", controllers.AddSegmentHandler)
	segments.Put("/", controllers.UpdateSegmentHandler)
	segments.Delete("/", controllers.DeleteSegmentHandler)

	segments.Patch("/add_user", controllers.AddSegmentsToUserHandler)
	segments.Patch("/remove_user", controllers.RemoveUserFromSegment)
	segments.Patch("/add_and_remove", controllers.AddAndRemoveSegmentsHandler)

	log.Fatal(app.Listen(":8000"))
}
