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
	db := initializers.DB
	services.Segments = services.CreateSegmentService(db)
	services.Users = services.CreateUserService(db)
	services.SegmentLogs = services.CreateSegmentLogService(db)
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/api/healthcheck", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Everything is fine.",
		})
	})

	app.Static("/static", "./static")

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

	segmentLogs := app.Group("/api/segment_logs")
	segmentLogs.Get("/", controllers.GetSegmentLogsHandler)

	log.Fatal(app.Listen(":8000"))
}
