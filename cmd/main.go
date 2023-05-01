package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
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

	// //Setting up routes
	router := routes.GetRouter()

	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST", "GET", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"www.example.com", "http://localhost:5500", "*"})
	fmt.Println("Listening on port 3000...")
	log.Fatal(http.ListenAndServe(os.Getenv("APP_PORT"), handlers.CORS(credentials, methods, origins)(router)))

}
