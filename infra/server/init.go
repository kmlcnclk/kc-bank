package server

import (
	errorresponse "kc-bank/pkg/error_response"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Init() *fiber.App {
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {

			// Default status code
			statusCode := fiber.StatusInternalServerError

			// Retrieve the custom error from fiber's context if it exists
			var customError errorresponse.CustomError

			if e, ok := err.(*fiber.Error); ok {
				// Fiber error, use its status code and message
				statusCode = e.Code
				customError = errorresponse.CustomError{
					StatusCode: statusCode,
					Message:    e.Message,
				}
			} else {

				// Non-fiber error, use default status code and message
				customError = errorresponse.CustomError{
					StatusCode: statusCode,
					Message:    err.Error(),
				}
			}

			// Send custom error response
			return ctx.Status(customError.StatusCode).JSON(customError)
		},
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Concurrency:  256 * 1024,
	})

	return app
}
