package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/remaster/webauthn/pkg/config"
)

type SessionData struct {
	Challenge            string    `json:"challenge"`
	UserID               uint64    `gorm:"primaryKey"`
	UserDisplayName      string    `json:"user_display_name"`
	allowedCredentialIDs [][]byte  `json:"allowed_credentials,omitempty", gorm:"-:migration, type:text"`
	Expires              time.Time `json:"expires"`

	UserVerification protocol.UserVerificationRequirement `json:"userVerification"`
	extensions       protocol.AuthenticationExtensions    `json:"extensions,omitempty", gorm:"-:migration, type:text"`
}

// Initialization function
func init() {
	config.ConnectDB()
	db = config.GetDB()
	db.AutoMigrate(&SessionData{})
}

func CreateSession(challenge string, userId uint64, displayName string, expires time.Time, userVerification protocol.UserVerificationRequirement) SessionData {
	rand.Seed(time.Now().UnixNano())
	session := &SessionData{}
	session.Challenge = challenge
	session.UserID = userId
	session.UserDisplayName = displayName
	session.Expires = expires
	session.UserVerification = userVerification
	result := db.Create(&session)
	if result.Error != nil {
		fmt.Print(result.Error)
	}
	return *session
}

func GetSessionByUserId(userId uint64) (SessionData, error) {
	var getSession SessionData
	result := db.Where("user_id=?", userId).First(&getSession)
	if result.Error != nil {
		return getSession, result.Error
	}
	return getSession, nil
}

func UpdateSessionByUserId(challenge string, userId uint64, displayName string, expires time.Time, userVerification protocol.UserVerificationRequirement) (SessionData, error) {
	var updateSession SessionData
	result := db.Model(&updateSession).Where("user_id=?", userId).Updates(SessionData{
		Challenge:        challenge,
		UserID:           userId,
		UserDisplayName:  displayName,
		Expires:          expires,
		UserVerification: userVerification,
	})
	if result.Error != nil {
		return updateSession, result.Error
	}
	return updateSession, nil
}
