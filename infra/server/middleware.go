package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func InitMiddlewares(app *fiber.App) {
	// Recover Middleware
	app.Use(recover.New())
}
