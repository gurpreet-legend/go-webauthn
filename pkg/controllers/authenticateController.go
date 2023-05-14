package controllers

// func BeginLogin(w http.ResponseWriter, r *http.Request) {

// 	//Get user by username
// 	params := mux.Vars(r)
// 	username := params["username"]
// 	user, err := models.GetUserByName(username)
// 	if err != nil {
// 		fmt.Println("Error while logging user.")
// 		fmt.Println(err)
// 		utils.JsonResponse(w, err, http.StatusInternalServerError)
// 		return
// 	}

// 	web := config.GetWebAuthn()

// 	// get credential for username
// 	creds, err := models.GetCredentialByUserId(user.Id)
// 	if err != nil {
// 		fmt.Println("Error while getting credentials for user.")
// 		fmt.Println(err)
// 		utils.JsonResponse(w, err, http.StatusInternalServerError)
// 		return
// 	}

// 	webCred := webauthn.Credential{
// 		ID:              creds[0].CredID,
// 		PublicKey:       creds[0].PublicKey,
// 		AttestationType: creds[0].AttestationType,
// 		Authenticator: webauthn.Authenticator{
// 			AAGUID:       creds[0].Authenticator.AAGUID,
// 			SignCount:    creds[0].Authenticator.SignCount,
// 			CloneWarning: creds[0].Authenticator.CloneWarning,
// 		},
// 	}
// 	user.AddCredential(webCred)
// 	options, sessionData, err := web.BeginLogin(user)
// 	if err != nil {
// 		fmt.Println("Error while BeginLogin function.")
// 		fmt.Println(err)
// 		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// store session data as marshaled JSON
// 	_, err = models.GetSessionByUserId(user.Id)
// 	if err != nil {
// 		models.CreateSession(sessionData.Challenge, user.Id, sessionData.UserVerification)
// 	} else {
// 		models.UpdateSessionByUserId(sessionData.Challenge, user.Id, sessionData.UserVerification)
// 	}

// 	utils.JsonResponse(w, options, http.StatusOK)
// }

// type CredentialRequest struct {
// 	Id       string `json:"id"`
// 	RawId    string `json:"rawId"`
// 	Type     string `json:"type"`
// 	Response struct {
// 		AuthenticatorData string `json:"authenticatorData"`
// 		ClientDataJSON    string `json:"clientDataJSON"`
// 		Signature         string `json:"signature"`
// 		UserHandle        string `json:"userHandle"`
// 	} `json:"response"`
// }

// func FinishLogin(w http.ResponseWriter, r *http.Request) {
// 	// parsing credential response body
// 	var responseBody CredentialRequest
// 	err := json.NewDecoder(r.Body).Decode(&responseBody)
// 	if err != nil {
// 		fmt.Println("Error while parsing response body")
// 		fmt.Println(err)
// 	}

// 	data, _ := json.Marshal(responseBody)
// 	reader := bytes.NewReader(data)
// 	response, _ := protocol.ParseCredentialRequestResponseBody(reader)

// 	//Get user by username
// 	params := mux.Vars(r)
// 	username := params["username"]
// 	user, err := models.GetUserByName(username)
// 	if err != nil {
// 		fmt.Println("Error while logging user.")
// 		fmt.Println(err)
// 		utils.JsonResponse(w, err, http.StatusInternalServerError)
// 		return
// 	}

// 	// get webauthn instance
// 	web := config.GetWebAuthn()

// 	// Load the session data
// 	sess, err := models.GetSessionByUserId(user.Id)
// 	if err != nil {
// 		fmt.Println("Error while getting session from db")
// 		fmt.Println(err)
// 		return
// 	}

// 	// Create instance of webauthn SessionData
// 	webSessionData := webauthn.SessionData{
// 		Challenge:        sess.Challenge,
// 		UserID:           utils.ConvertIntToByteArray(sess.UserID),
// 		UserVerification: sess.UserVerification,
// 	}

// 	// get credential for username
// 	creds, err := models.GetCredentialByUserId(user.Id)
// 	if err != nil {
// 		fmt.Println("Error while getting credentials for user.")
// 		fmt.Println(err)
// 		utils.JsonResponse(w, err, http.StatusInternalServerError)
// 		return
// 	}
// 	webCred := webauthn.Credential{
// 		ID:              creds[0].CredID,
// 		PublicKey:       creds[0].PublicKey,
// 		AttestationType: creds[0].AttestationType,
// 		Authenticator: webauthn.Authenticator{
// 			AAGUID:       creds[0].Authenticator.AAGUID,
// 			SignCount:    creds[0].Authenticator.SignCount,
// 			CloneWarning: creds[0].Authenticator.CloneWarning,
// 		},
// 	}
// 	user.AddCredential(webCred)

// 	_, err = web.ValidateLogin(user, webSessionData, response)
// 	if err != nil {
// 		fmt.Println("Error while validating credentials.")
// 		fmt.Println(err)
// 		utils.JsonResponse(w, err, http.StatusInternalServerError)
// 		return
// 	}
// 	fmt.Println("LOGIN SUCCESSFULLLLL -----------------")

// 	// handle successful login
// 	utils.JsonResponse(w, "Login Success", http.StatusOK)
// }
