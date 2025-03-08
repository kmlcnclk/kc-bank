package main

import (
	"go.uber.org/zap"

	"kc-bank/app/controllers/healthcheck"
	userController "kc-bank/app/controllers/user"
	"kc-bank/app/repository"
	"kc-bank/app/services/command"
	"kc-bank/app/services/query"
	"kc-bank/infra/couchbase"
	"kc-bank/infra/server"
	"kc-bank/pkg/config"
	_ "kc-bank/pkg/log"
	"kc-bank/pkg/services"
)

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()

	zap.L().Info("app starting...")

	// rmq, err := rabbitmq.NewRabbitMQ(appConfig.RabbitMQURL, "my_queue", "my_exchange", "direct")

	// if err != nil {
	// 	zap.L().Fatal("failed to initialize RabbitMQ", zap.Error(err))
	// }

	// defer rmq.Close()

	// Initialize Couchbase
	cb, err := couchbase.NewCouchbase(appConfig.CouchbaseUrl, appConfig.CouchbaseUsername, appConfig.CouchbasePassword)

	if err != nil {
		zap.L().Fatal("Failed to initialize Couchbase instance", zap.Error(err))
	}

	// Initialize user bucket
	userBucket := cb.InitializeBucket("users")

	// Dependency Injection
	userRepository := repository.NewUserRepository(userBucket)
	passwordService := services.NewPasswordService()
	userCommand := command.NewCommandHandler(userRepository, passwordService)
	userQuery := query.NewUserQueryService(userRepository)

	// Initialize controllers
	getUserHandler := userController.NewGetUserHandler(userQuery)
	createUserHandler := userController.NewCreateUserHandler(userCommand)

	// Initialize healthcheck handler
	healthcheckHandler := healthcheck.NewHealthCheckHandler()

	// Init Fiber app
	app := server.Init()

	// Init middlewares
	server.InitMiddlewares(app)

	// Init routers
	server.InitRouters(app, getUserHandler, createUserHandler, healthcheckHandler)

	// Start server
	server.Start(app, appConfig)

	// Graceful shutdown
	server.GracefulShutdown(app)
}
