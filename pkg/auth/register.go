package auth

import (
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
)

func BeginRegistration(w http.ResponseWriter, r *http.Request) {
	user := datastore.GetUser() // Find or create the new user
	options, session, err := web.BeginRegistration(user)
	// handle errors if present
	// store the sessionData values
	JSONResponse(w, options, http.StatusOK) // return the options generated
	// options.publicKey contain our registration options
}

func FinishRegistration(w http.ResponseWriter, r *http.Request) {
	response, err := protocol.ParseCredentialCreationResponseBody(r.Body)
	if err != nil {
		// Handle Error and return.

		return
	}

	user := datastore.GetUser() // Get the user

	// Get the session data stored from the function above
	session := datastore.GetSession()

	credential, err := w.CreateCredential(user, session, response)
	if err != nil {
		// Handle Error and return.

		return
	}

	// If creation was successful, store the credential object
	JSONResponse(w, "Registration Success", http.StatusOK) // Handle next steps

	// Pseudocode to add the user credential.
	user.AddCredential(credential)
	datastore.SaveUser(user)
}
