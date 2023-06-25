package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/mux"
	"github.com/remaster/webauthn/pkg/config"
	"github.com/remaster/webauthn/pkg/models"
	"github.com/remaster/webauthn/pkg/utils"
)

var ErrCredentialsInvalid = errors.New("Error while validating credentials.")

func BeginLogin(w http.ResponseWriter, r *http.Request) {

	//Get user by username
	params := mux.Vars(r)
	username := params["username"]
	user, err := models.GetUserByName(username)
	if err != nil {
		if err == models.ErrUserNotFound {
			errorMessage := fmt.Sprintf("User '%s' not found.", username)
			fmt.Println(errorMessage)
			utils.JsonResponse(w, errorMessage, http.StatusNotFound)
			return
		}

		// Other error occurred
		fmt.Println(err)
		utils.JsonResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(user.Credentials) == 0 {
		totpFallback()
		utils.JsonResponse(w, err, http.StatusUnauthorized)
		return
	}

	web := config.GetWebAuthn()

	options, sessionData, err := web.BeginLogin(user)
	if err != nil {
		fmt.Println("Error while BeginLogin function.")
		fmt.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("OPTIONS------\n %+v\n", options)
	fmt.Printf("SESSIONDATA------\n %+v\n", sessionData)

	// store session data as marshaled JSON
	_, err = models.GetSessionByUserId(user.Id)
	if err != nil {
		models.CreateSession(sessionData)
	} else {
		models.UpdateSessionByUserId(user.Id, sessionData)
	}

	utils.JsonResponse(w, options, http.StatusOK)
}

type CredentialRequest struct {
	Id       string `json:"id"`
	RawId    string `json:"rawId"`
	Type     string `json:"type"`
	Response struct {
		AuthenticatorData string `json:"authenticatorData"`
		ClientDataJSON    string `json:"clientDataJSON"`
		Signature         string `json:"signature"`
		UserHandle        string `json:"userHandle"`
	} `json:"response"`
}

func FinishLogin(w http.ResponseWriter, r *http.Request) {
	// parsing credential response body
	var responseBody CredentialRequest
	err := json.NewDecoder(r.Body).Decode(&responseBody)
	if err != nil {
		fmt.Println("Error while parsing response body")
		fmt.Println(err)
	}

	data, _ := json.Marshal(responseBody)
	reader := bytes.NewReader(data)
	response, _ := protocol.ParseCredentialRequestResponseBody(reader)

	fmt.Printf("PARSED CREDENTIAL RESPONSE-----\n%+v\n", response)
	//Get user by username
	params := mux.Vars(r)
	username := params["username"]
	user, err := models.GetUserByName(username)
	if err != nil {
		fmt.Println("Error while logging user.")
		fmt.Println(err)
		utils.JsonResponse(w, err, http.StatusInternalServerError)
		return
	}
	fmt.Printf("USER HANDLE-----\n%v\n", utils.ConvertByteArrayToInt(response.Response.UserHandle))
	fmt.Printf("USER HANDLE-----\n%v\n", utils.ConvertByteArrayToInt(user.WebAuthnID()))

	// get webauthn instance
	web := config.GetWebAuthn()

	// Load the session data
	sess, err := models.GetSessionByUserId(user.Id)
	if err != nil {
		fmt.Println("Error while getting session from db")
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

	_, err = web.ValidateLogin(user, webSessionData, response)
	if err != nil {
		fmt.Println(user.Credentials)
		fmt.Println("Error while validating credentials.")
		fmt.Println(err)
		utils.JsonResponse(w, err, http.StatusUnauthorized)
		return
	}

	// handle successful login
	utils.JsonResponse(w, "Login Success", http.StatusOK)
}
