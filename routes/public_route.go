package routes

import (
	"github.com/gofiber/fiber/v3"
)

func Public(app *fiber.App) {
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// publicV1Endpoint := app.Group("/api/public/v1")

	// monitorEndpoint := publicV1Endpoint.Group("/monitor")
	// monitorEndpoint.Get("/health", ctrlMonitor.HealthCheckStatus)
}
