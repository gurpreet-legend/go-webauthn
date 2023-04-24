package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/remaster/webauthn/pkg/config"
	"github.com/remaster/webauthn/pkg/models"
)

func beginRegistration(c *fiber.Ctx) error {
	user := models.CreateUser() // Create the new user
	web := config.GetWebAuthn()
	options, session, err := web.BeginRegistration(*user)
	if err != nil {
		fmt.Print("BeginRegistration failed!!")
		c.Status(400).JSON(&fiber.Map{
			"error": err,
		})
		return err
	}
	// store the sessionData values
	fmt.Print(session)
	// ---> still need to implement the store function

	// return the options generated
	c.Status(200).JSON(&fiber.Map{
		"options": options,
	})
	return nil
	// options.publicKey contain our registration options
}

// func FinishRegistration(c *fiber.Ctx) error {
// 	response, err := protocol.ParseCredentialCreationResponseBody()
// 	if err != nil {
// 		fmt.Print("BeginRegistration failed!!")
// 		c.Status(400).JSON(&fiber.Map{
// 			"error": err,
// 		})
// 		return err
// 	}

// 	user := models.GetUserById() // Get the user

// 	// Get the session data stored from the function above
// 	session := datastore.GetSession()

// 	credential, err := w.CreateCredential(user, session, response)
// 	if err != nil {
// 		// Handle Error and return.

// 		return
// 	}

// 	// If creation was successful, store the credential object
// 	JSONResponse(w, "Registration Success", http.StatusOK) // Handle next steps

// 	// Pseudocode to add the user credential.
// 	user.AddCredential(credential)
// 	datastore.SaveUser(user)
// }

func RegisterController(c *fiber.Ctx) error {
	beginRegistration(c)
	return nil
}
