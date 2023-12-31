package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/mux"
	"github.com/remaster/webauthn/pkg/config"
	"github.com/remaster/webauthn/pkg/models"
	"github.com/remaster/webauthn/pkg/utils"
)

func BeginRegistration(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	//Get user by username
	user, err := models.GetUserByName(username)
	if err != nil {
		displayName := strings.Split(username, "@")[0]
		models.CreateUser(username, displayName)
		user, _ = models.GetUserByName(username)
	}

	web := config.GetWebAuthn()

	fmt.Printf("BEGIN REGISTER USERID-----\n%v\n", user.Id)
	// Begin registration using user
	options, sessionData, err := web.BeginRegistration(user)
	if err != nil {
		fmt.Println("Error while begin registartion function")
		fmt.Println(err)
		return
	}

	// Storing session
	models.CreateSession(sessionData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	utils.JsonResponse(w, options, http.StatusOK)
}

type AttestationResponse struct {
	Id       string `json:"id"`
	RawId    string `json:"rawId"`
	Type     string `json:"type"`
	Response struct {
		AttestationObject string `json:"attestationObject"`
		ClientDataJSON    string `json:"clientDataJSON"`
	} `json:"response"`
}

func FinishRegistration(w http.ResponseWriter, r *http.Request) {
	var responseBody AttestationResponse
	err := json.NewDecoder(r.Body).Decode(&responseBody)
	if err != nil {
		fmt.Println("Error while parsing response body")
		fmt.Println(err)
	}

	// Parsing data and creating a io.Reader for response body
	data, _ := json.Marshal(responseBody)
	reader := bytes.NewReader(data)

	//Get user by username
	params := mux.Vars(r)
	username := params["username"]
	user, err := models.GetUserByName(username)
	if err != nil {
		fmt.Println("Error while getting user from db")
		fmt.Println(err)
		return
	}

	fmt.Printf("FINISH REGISTER USERID-----\n%v\n", user.Id)

	// Load the session data
	sess, err := models.GetSessionByUserId(user.Id)
	if err != nil {
		fmt.Println("Error while getting session from db")
		fmt.Println(err)
		return
	}

	response, err := protocol.ParseCredentialCreationResponseBody(reader)
	if err != nil {
		fmt.Println("Error while parsing credential creation response body")
		fmt.Println(err)
		return
	}

	// Create instance of webauthn SessionData
	webSessionData := webauthn.SessionData{
		Challenge:            sess.Challenge,
		UserID:               utils.ConvertIntToByteArray(sess.UserID),
		UserVerification:     sess.UserVerification,
		Extensions:           sess.Extensions,
		AllowedCredentialIDs: sess.AllowedCredentialIDs,
	}

	// Get instance of webauthn
	web := config.GetWebAuthn()

	// Create credentials for user
	credential, err := web.CreateCredential(user, webSessionData, response)
	if err != nil {
		fmt.Println("Error while creating credentials")
		fmt.Println(err)
	}

	// Storing credentials
	_, err = models.AddCredentialToUser(user.Id, credential)
	if err != nil {
		fmt.Println("ADD CREDENTIAL ERROR")
		fmt.Println(err)
	}

	utils.JsonResponse(w, "Registration success", http.StatusOK)
}
