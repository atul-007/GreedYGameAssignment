package routes

import (
	handler "github.com/atul-007/GreedyGameAssignment/handlers"
	"github.com/atul-007/GreedyGameAssignment/services"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App, appHandlers handler.Handlers, allServices services.Services) {
	api := app.Group("/api")
	dataApi := api.Group("/data")

	InitDataRoutes(dataApi, appHandlers)

}
func InitDataRoutes(api fiber.Router, appHandlers handler.Handlers) {
	api.Post("/", appHandlers.Data.HandleDb)
}
