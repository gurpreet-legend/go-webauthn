package auth

func BeginLogin(w http.ResponseWriter, r *http.Request) {
	user := datastore.GetUser() // Find the user

	options, session, err := w.BeginLogin(user)
	if err != nil {
		// Handle Error and return.

		return
	}

	// store the session values
	datastore.SaveSession(session)

	JSONResponse(w, options, http.StatusOK) // return the options generated
	// options.publicKey contain our registration options
}

func FinishLogin(w http.ResponseWriter, r *http.Request) {
	response, err := protocol.ParseCredentialRequestResponseBody(r.Body)
	if err != nil {
		// Handle Error and return.

		return
	}

	user := datastore.GetUser() // Get the user

	// Get the session data stored from the function above
	session := datastore.GetSession()

	credential, err := w.ValidateLogin(user, session, response)
	if err != nil {
		// Handle Error and return.

		return
	}

	// If login was successful, handle next steps
	JSONResponse(w, "Login Success", http.StatusOK)
}
