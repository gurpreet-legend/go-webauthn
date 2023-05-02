package controllers

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gorilla/mux"
	"github.com/remaster/webauthn/pkg/config"
	"github.com/remaster/webauthn/pkg/models"
	"github.com/remaster/webauthn/pkg/utils"
)

func BeginRegistration(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	fmt.Printf("%v\n", username)

	//Get user by username
	user, err := models.GetUserByName(username)
	if err != nil { // Create user if not registered already
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

	fmt.Println(uint64(binary.LittleEndian.Uint64(sessionData.UserID)))

	// Storing session
	_, err = models.GetSessionByUserId(user.Id)
	if err != nil {
		models.CreateSession(sessionData.Challenge, user.Id, user.DisplayName, sessionData.Expires, sessionData.UserVerification)
	} else {
		models.UpdateSessionByUserId(sessionData.Challenge, user.Id, user.DisplayName, sessionData.Expires, sessionData.UserVerification)
	}

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
	fmt.Println("RESPONSE BODY---------------")
	fmt.Printf("%+v\n", responseBody)

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
	fmt.Println("USER DETAILS---------------")
	fmt.Println(user)

	// Load the session data
	sess, err := models.GetSessionByUserId(user.Id)
	if err != nil {
		fmt.Println("Error while getting session from db")
		fmt.Println(err)
		return
	}
	fmt.Println("SESSION DETAILS---------------")
	fmt.Printf("%+v\n", sess)

	// testing .....
	// var ccr protocol.CredentialCreationResponse
	// json.NewDecoder(reader).Decode(&ccr)
	// fmt.Printf("%+v", ccr)

	// var ccr protocol.CredentialCreationResponse

	// if err = json.NewDecoder(reader).Decode(&ccr); err != nil {
	// 	return nil
	// }
	// ccrAttest := ccr.AttestationResponse
	// p := &protocol.ParsedAttestationResponse{}
	// if err = json.Unmarshal(ccrAttest.ClientDataJSON, &p.CollectedClientData); err != nil {
	// 	return nil
	// }

	// fmt.Println("CCR CREDENTIALS START---------------")
	// // ccr1, _ := ccr.AttestationResponse.Parse()
	// fmt.Printf("%+v", p.CollectedClientData)
	// fmt.Println("CCR CREDENTIALS END---------------")

	response, err := protocol.ParseCredentialCreationResponseBody(reader)
	if err != nil {
		fmt.Println("Error while parsing credential creation response body")
		fmt.Println(err)
		return
	}
	fmt.Println("PARSED CREDENTIALS---------------")
	fmt.Printf("%+v\n", response)

	// Create instance of webauthn SessionData
	webSessionData := webauthn.SessionData{
		Challenge:        sess.Challenge,
		UserID:           utils.ConvertIntToByteArray(sess.UserID),
		UserDisplayName:  sess.UserDisplayName,
		Expires:          sess.Expires,
		UserVerification: sess.UserVerification,
	}

	// Get instance of webauthn
	web := config.GetWebAuthn()

	// Create credentials for user
	credential, err := web.CreateCredential(user, webSessionData, response)
	if err != nil {
		fmt.Println("Error while creating credentials")
		fmt.Println(err)
	}

	// data, _ := json.Marshal(responseBody)
	// reader := bytes.NewReader(data)

	// req, _ := http.NewRequest("FinishRegistrationRequest", "http://localhost:3000/", reader)

	// credential, err := web.FinishRegistration(user, webSessionData, req)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	fmt.Println("CREDENTIAL DETAILS---------------")
	fmt.Printf("%+v\n", credential)
}
