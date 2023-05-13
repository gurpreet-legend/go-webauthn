package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"github.com/remaster/webauthn/pkg/config"
	"github.com/remaster/webauthn/pkg/models"
	"github.com/remaster/webauthn/pkg/routes"
)

// Your initialization function
func main() {
	//dotenv
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//WebAuthn setup
	config.SetupWebAuthn()

	// //Setting up routes
	router := routes.GetRouter()

	// Initialize couchbase
	config.InitCouchbaseConnection()

	err = models.CouchCreateUser("wiz3@singh", "wiz3")
	fmt.Println("HERE COMES THE ERROR")
	fmt.Println(err)

	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST", "GET", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:5500", "*"})
	fmt.Println("Listening on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(credentials, methods, origins)(router)))

}
