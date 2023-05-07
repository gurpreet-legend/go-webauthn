package controllers

import (
	"fmt"
	"net/http"

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
		return
	}

	web := config.GetWebAuthn()

	// generate PublicKeyCredentialRequestOptions, session data
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
