package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/gofiber/fiber/v2"
	"github.com/remaster/webauthn/pkg/config"
	"github.com/remaster/webauthn/pkg/models"
)

func BeginRegistration(c *fiber.Ctx) error {
	username := c.Params("username")
	fmt.Printf("%v\n", username)

	//Get user by username
	user, err := models.GetUserByName(username)
	if err != nil { // Create user if not registered already
		displayName := strings.Split(username, "@")[0]
		fmt.Println(displayName)
		models.CreateUser(username, displayName)
	}

	web := config.GetWebAuthn()

	// Begin registration using user
	options, sessionData, err := web.BeginRegistration(user)
	if err != nil {
		fmt.Println(err)
		c.Status(500).JSON(&fiber.Map{
			"error":   err,
			"message": "Internal server error. Registration failed \n",
		})
		return err
	}

	// Storing session
	sess := *config.GetSession()
	sessData := *sessionData
	fmt.Println(sessData.UserID)
	store, err := sess.Get(c)
	if err != nil {
		fmt.Println(err)
		c.Status(500).JSON(&fiber.Map{
			"error":   err,
			"message": "Internal server error. Session storage failed \n",
		})
		return err
	}
	store.Set("registration", sessData)
	store.Save()

	c.Status(200).JSON(&fiber.Map{
		"options": options,
	})

	return nil
}

type NestedResponse struct {
	AttestationObject string `json:"attestationObject"`
	ClientDataJSON    string `json:"clientDataJSON"`
}
type AttestationResponse struct {
	Id       string         `json:"id"`
	RawId    string         `json:"rawId"`
	Type     string         `json:"type"`
	Response NestedResponse `json:"response"`
}

func FinishRegistration(c *fiber.Ctx) error {
	var responseBody AttestationResponse
	if err := c.BodyParser(&responseBody); err != nil {
		c.Status(500).JSON(&fiber.Map{
			"error":   err,
			"message": "Internal server error. Error while parsing the body",
		})
		fmt.Println("Internal server error. Error while parsing the body")
	}

	// Parsing data and creating a io.Reader for response body
	data, _ := json.Marshal(responseBody)
	reader := bytes.NewReader(data)

	fmt.Printf("\n\n\n")

	_, err := protocol.ParseCredentialCreationResponseBody(reader)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Final Registration failed!!")
		c.Status(400).JSON(&fiber.Map{
			"error": err,
		})
		return err
	}
	// fmt.Println(response)

	username := c.Params("username")

	//Get user by username
	_, err = models.GetUserByName(username)
	if err != nil { // Create user if not registered already
		c.Status(500).JSON(&fiber.Map{
			"error":   err,
			"message": "Internal server error. User not found during validating attestation. Please register the user first.",
		})
	}

	// load the session data
	sess := *config.GetSession()
	store, err := sess.Get(c)
	if err != nil {
		c.Status(500).JSON(&fiber.Map{
			"error":   err,
			"message": "Internal server error. Unable to get the session storage.",
		})
	}
	sessionData := store.Get("registration")
	fmt.Println(sessionData)

	// web := config.GetWebAuthn()

	// credential, err := web.FinishRegistration(user, sessionData, c.Response())

	// credential, err := web.CreateCredential(user, sessionData, response)
	// if err != nil {
	// 	fmt.Println(err)
	// 	c.Status(500).JSON(&fiber.Map{
	// 		"error":   err,
	// 		"message": "Internal server error. Unable to generate credentials.",
	// 	})
	// 	return
	// }

	// user.AddCredential(*credential)
	return nil
}

func RegisterController(c *fiber.Ctx) error {
	// beginRegistration(c)
	return nil
}
