package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	// config.CreateSession()

	app := fiber.New()

	//cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5500",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	//Setting up routes
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
