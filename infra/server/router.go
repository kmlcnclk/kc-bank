package server

import (
	"kc-bank/app/controllers/healthcheck"
	"kc-bank/app/controllers/user"
	"kc-bank/pkg/handler"

	"github.com/gofiber/fiber/v2"
)

func InitRouters(app *fiber.App,
	getUserHandler *user.GetUserHandler,
	createUserHandler *user.CreateUserHandler,
	healthcheckHandler *healthcheck.HealthCheckHandler,

) {

	app.Get("/healthcheck", handler.Handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheckHandler))

	userGroup := app.Group("/api/v1/user")

	userGroup.Get("/:id", handler.Handle[user.GetUserRequest, user.GetUserResponse](getUserHandler))
	userGroup.Post("/", handler.Handle[user.CreateUserRequest, user.CreateUserResponse](createUserHandler))
}
