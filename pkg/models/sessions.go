package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/remaster/webauthn/pkg/config"
	"github.com/remaster/webauthn/pkg/utils"
)

type SessionData struct {
	Challenge string `json:"challenge"`
	UserID    uint64 `json:"userId"`
	// UserDisplayName      string    `json:"user_display_name"`
	AllowedCredentialIDs [][]byte `json:"allowedCredentials"`
	// Expires              time.Time `json:"expires"`

	UserVerification protocol.UserVerificationRequirement `json:"userVerification"`
	Extensions       protocol.AuthenticationExtensions    `json:"extensions"`
}

func CreateSession(sessionData *webauthn.SessionData) error {
	session := &SessionData{
		Challenge:            sessionData.Challenge,
		UserID:               utils.ConvertByteArrayToInt(sessionData.UserID),
		UserVerification:     sessionData.UserVerification,
		AllowedCredentialIDs: sessionData.AllowedCredentialIDs,
		Extensions:           sessionData.Extensions,
	}

	sessions := config.GetDefaultScope().Collection("sessions")
	_, err := sessions.Upsert(strconv.FormatUint(session.UserID, 10), session, &gocb.UpsertOptions{
		Timeout: 5 * time.Second,
		Expiry:  10 * time.Minute,
	})
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error creating session.")
		return err
	}
	return nil
}

func GetSessionByUserId(userId uint64) (SessionData, error) {
	scope := config.GetDefaultScope()
	var getSession SessionData
	queryString := "SELECT x.* FROM `webauthn-bucket`.`webauthn-scope`.`sessions` x WHERE META().id=$userId"
	result, dbErr := config.ExecuteDBQuery(scope, queryString, &config.DBQueryParameters{
		"userId": strconv.FormatUint(userId, 10),
	})
	if dbErr != nil {
		fmt.Println(dbErr)
		fmt.Println("Error getting session by userId.")
		return getSession, dbErr
	}
	err := result.One(&getSession)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error getting session by userId.")
		return getSession, err
	}
	return getSession, nil
}

func UpdateSessionByUserId(userId uint64, sessionData *webauthn.SessionData) (SessionData, error) {
	var updateSession SessionData

	session := &SessionData{
		Challenge:            sessionData.Challenge,
		UserID:               utils.ConvertByteArrayToInt(sessionData.UserID),
		UserVerification:     sessionData.UserVerification,
		AllowedCredentialIDs: sessionData.AllowedCredentialIDs,
		Extensions:           sessionData.Extensions,
	}

	sessions := config.GetDefaultScope().Collection("sessions")
	_, err := sessions.Replace(strconv.FormatUint(session.UserID, 10), session, &gocb.ReplaceOptions{
		Timeout: 5 * time.Second,
		Expiry:  10 * time.Minute,
	})

	if err != nil {
		fmt.Println(err)
		fmt.Println("Error creating session.")
		return *session, err
	}
	return updateSession, nil
}
