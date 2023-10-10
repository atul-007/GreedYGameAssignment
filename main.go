package main

import (
	"log"

	handler "github.com/atul-007/GreedyGameAssignment/handlers"
	"github.com/atul-007/GreedyGameAssignment/routes"
	"github.com/atul-007/GreedyGameAssignment/services"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
	})

	// Cors
	//app.Use(cors.New())
	services := services.Init()
	handlers := handler.Init()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Init routes
	routes.Init(app, handlers, services)

	log.Fatal(app.Listen(":3000"))
}
