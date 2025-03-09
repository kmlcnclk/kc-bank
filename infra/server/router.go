package server

import (
	"kc-bank/app/controllers/account"
	"kc-bank/app/controllers/healthcheck"
	"kc-bank/app/controllers/user"
	"kc-bank/pkg/handler"

	"github.com/gofiber/fiber/v2"
)

func InitRouters(app *fiber.App,
	getUserHandler *user.GetUserHandler,
	createUserHandler *user.CreateUserHandler,
	getUserAllHandler *user.GetUserAllHandler,
	healthcheckHandler *healthcheck.HealthCheckHandler,
	getAccountHandler *account.GetAccountHandler,
	getAccountAllHandler *account.GetAccountAllHandler,
	createAccountHandler *account.CreateAccountHandler,
	transferMoneyHandler *account.TransferMoneyHandler,
	transferMoneyWithRabbitMQHandler *account.TransferMoneyWithRabbitMQHandler,
) {

	app.Get("/healthcheck", handler.Handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheckHandler))

	// User
	userGroup := app.Group("/api/v1/user")

	userGroup.Get("/", handler.Handle[user.GetUserAllRequest, user.GetUserAllResponse](getUserAllHandler))
	userGroup.Get("/:id", handler.Handle[user.GetUserRequest, user.GetUserResponse](getUserHandler))
	userGroup.Post("/", handler.Handle[user.CreateUserRequest, user.CreateUserResponse](createUserHandler))

	// Account
	accountGroup := app.Group("/api/v1/account")

	accountGroup.Get("/", handler.Handle[account.GetAccountAllRequest, account.GetAccountAllResponse](getAccountAllHandler))
	accountGroup.Get("/:id", handler.Handle[account.GetAccountRequest, account.GetAccountResponse](getAccountHandler))
	accountGroup.Post("/", handler.Handle[account.CreateAccountRequest, account.CreateAccountResponse](createAccountHandler))
	accountGroup.Post("/transfer-money", handler.Handle[account.TransferMoneyRequest, account.TransferMoneyResponse](transferMoneyHandler))
	accountGroup.Post("/transfer-money-with-rmq", handler.Handle[account.TransferMoneyWithRabbitMQRequest, account.TransferMoneyWithRabbitMQResponse](transferMoneyWithRabbitMQHandler))
}
