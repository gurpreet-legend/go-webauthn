package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/remaster/webauthn/pkg/config"
	"github.com/remaster/webauthn/pkg/routes"
)

// Your initialization function
func main() {
	//dotenv
	godotenv.Load()

	//WebAuthn setup
	config.SetupWebAuthn()

	//Setup session
	config.CreateSession()

	//Setting up routes
	app := fiber.New()
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
