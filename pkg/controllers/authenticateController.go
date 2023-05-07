package controllers

import (
	"fmt"
	"net/http"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/mux"
	"github.com/remaster/webauthn/pkg/config"
	"github.com/remaster/webauthn/pkg/models"
	"github.com/remaster/webauthn/pkg/utils"
)

func BeginLogin(w http.ResponseWriter, r *http.Request) {

	//Get user by username
	params := mux.Vars(r)
	username := params["username"]
	user, err := models.GetUserByName(username)
	if err != nil {
		fmt.Println("Error while logging user, User doesn't exist.")
		fmt.Println(err)
		utils.JsonResponse(w, err, http.StatusInternalServerError)
		return
	}

	web := config.GetWebAuthn()

	// get credential for username
	creds, err := models.GetCredentialByUserId(user.Id)
	if err != nil {
		fmt.Println("Error while getting credentials for user.")
		fmt.Println(err)
		utils.JsonResponse(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Println("CREDENTIALS  FOR USERID ---------------")
	fmt.Printf("%+v\n", creds)

	webCred := webauthn.Credential{
		ID:              creds[0].CredID,
		PublicKey:       creds[0].PublicKey,
		AttestationType: creds[0].AttestationType,
		Authenticator: webauthn.Authenticator{
			AAGUID:       creds[0].Authenticator.AAGUID,
			SignCount:    creds[0].Authenticator.SignCount,
			CloneWarning: creds[0].Authenticator.CloneWarning,
		},
	}
	user.AddCredential(webCred)
	fmt.Println("USER WITH CREDENTIALS---------------")
	fmt.Printf("%+v", user)
	options, sessionData, err := web.BeginLogin(user)
	if err != nil {
		fmt.Println("Error while BeginLogin function.")
		fmt.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// store session data as marshaled JSON
	_, err = models.GetSessionByUserId(user.Id)
	if err != nil {
		models.CreateSession(sessionData.Challenge, user.Id, sessionData.UserVerification)
	} else {
		models.UpdateSessionByUserId(sessionData.Challenge, user.Id, sessionData.UserVerification)
	}

	utils.JsonResponse(w, options, http.StatusOK)
}
