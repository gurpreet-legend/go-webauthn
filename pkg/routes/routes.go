package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/remaster/webauthn/pkg/controllers"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/authenticate", controllers.AuthenticateController)
	app.Get("/register", controllers.RegisterController)
}
