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
	fmt.Println("BEGIN REGISTRATION STARTS--------------------")
	params := mux.Vars(r)
	username := params["username"]
	fmt.Printf("%v\n", username)

	//Get user by username
	user, err := models.GetUserByName(username)
	if err != nil {
		displayName := strings.Split(username, "@")[0]
		models.CreateUser(username, displayName)
		user, _ = models.GetUserByName(username)
	}

	web := config.GetWebAuthn()

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
	fmt.Println("BEGIN REGISTRATION END--------------------")
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
	fmt.Println("FINISH REGISTRATION STARTS--------------------")
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
	// credAuth := models.Authenticator{
	// 	AAGUID:       credential.Authenticator.AAGUID,
	// 	SignCount:    credential.Authenticator.SignCount,
	// 	CloneWarning: credential.Authenticator.CloneWarning,
	// }

	// _, err = models.GetCredentialByCredId(credential.ID)
	// if err != nil {
	// 	models.CreateCredential(user.Id, credential.ID, credential.PublicKey, credential.AttestationType, credAuth)
	// } else {
	// 	models.UpdateCredentialByUserId(user.Id, credential.ID, credential.PublicKey, credential.AttestationType, credAuth)
	// }

	utils.JsonResponse(w, "Registration success", http.StatusOK)
	fmt.Println("FINISH REGISTRATION ENDS--------------------")
}
