package server

import (
	"fmt"
	"kc-bank/pkg/config"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Start(app *fiber.App, appConfig *config.AppConfig) {

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", appConfig.Port)); err != nil {
			zap.L().Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	zap.L().Info("Server started on port", zap.String("port", appConfig.Port))

}
