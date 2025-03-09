package main

import (
	"go.uber.org/zap"

	accountController "kc-bank/app/controllers/account"
	"kc-bank/app/controllers/healthcheck"
	userController "kc-bank/app/controllers/user"
	"kc-bank/app/repository"
	accountCommand "kc-bank/app/services/account/command"
	accountQuery "kc-bank/app/services/account/query"
	userCommand "kc-bank/app/services/user/command"
	userQuery "kc-bank/app/services/user/query"
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

	cluster, err := couchbase.ConnectCouchbase(appConfig.CouchbaseUrl, appConfig.CouchbaseUsername, appConfig.CouchbasePassword)

	if err != nil {
		zap.L().Error("Failed to initialize Couchbase:", zap.Error(err))
	}

	cb := couchbase.NewCouchbase(cluster)

	// Initialize user bucket
	userBucket := cb.InitializeBucket("users")

	// Initialize user bucket
	accountBucket := cb.InitializeBucket("accounts")

	// Dependency Injection for User
	userRepository := repository.NewUserRepository(cluster, userBucket)
	passwordService := services.NewPasswordService()
	userCommand := userCommand.NewCommandHandler(userRepository, passwordService)
	userQuery := userQuery.NewUserQueryService(userRepository)

	// Dependency Injection for Account
	accountRepository := repository.NewAccountRepository(cluster, accountBucket)
	ibanService := services.NewIbanService()
	accountCommand := accountCommand.NewCommandHandler(accountRepository, ibanService)
	accountQuery := accountQuery.NewAccountQueryService(accountRepository)

	// Initialize controllers for User
	getUserHandler := userController.NewGetUserHandler(userQuery)
	getUserAllHandler := userController.NewGetUserAllHandler(userQuery)
	createUserHandler := userController.NewCreateUserHandler(userCommand)

	// Initialize controllers for Account
	getAccountHandler := accountController.NewGetAccountHandler(accountQuery)
	getAccountAllHandler := accountController.NewGetAccountAllHandler(accountQuery)
	createAccountHandler := accountController.NewCreateAccountHandler(accountCommand)
	transferMoneyHandler := accountController.NewTransferMoneyHandler(accountCommand)

	// Initialize healthcheck handler
	healthcheckHandler := healthcheck.NewHealthCheckHandler()

	// Init Fiber app
	app := server.Init()

	// Init middlewares
	server.InitMiddlewares(app)

	// Init routers
	server.InitRouters(
		app,
		getUserHandler,
		createUserHandler,
		getUserAllHandler,
		healthcheckHandler,
		getAccountHandler,
		getAccountAllHandler,
		createAccountHandler,
		transferMoneyHandler,
	)

	// Start server
	server.Start(app, appConfig)

	// Graceful shutdown
	server.GracefulShutdown(app)
}
