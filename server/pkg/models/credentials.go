package models

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/remaster/webauthn/pkg/config"
)

type Authenticator struct {
	AAGUID       []byte `json:"aaguId"`
	SignCount    uint32 `json:"signCount"`
	CloneWarning bool   `json:"cloneWarning"`
}
type Credential struct {
	ID              []byte        `json:"id"`
	PublicKey       []byte        `json:"publicKey"`
	AttestationType string        `json:"attestationType"`
	Authenticator   Authenticator `json:"authenticator"`
}

func AddCredentialToUser(userId uint64, credential *webauthn.Credential) (Credential, error) {
	rand.Seed(time.Now().UnixNano())
	cred := Credential{
		ID:              credential.ID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		Authenticator: Authenticator{
			AAGUID:       credential.Authenticator.AAGUID,
			SignCount:    credential.Authenticator.SignCount,
			CloneWarning: credential.Authenticator.CloneWarning,
		},
	}
	// fmt.Println("HERE IS THE CREDENTIAL-------------")
	// fmt.Printf("%+v\n", cred)
	// fmt.Println("USER ID TO STORE CRED-------------")
	// fmt.Printf("%+v\n", userId)
	scope := config.GetDefaultScope()
	fmt.Printf("%+v\n", cred)
	queryString := "UPDATE `webauthn-bucket`.`webauthn-scope`.`users` u SET u.credentials = ARRAY_DISTINCT(ARRAY_APPEND(u.credentials, $cred)) WHERE u.id=$userId"
	_, dbErr := config.ExecuteDBQuery(scope, queryString, &config.DBQueryParameters{
		"cred":   cred,
		"userId": strconv.FormatUint(userId, 10),
	})
	if dbErr != nil {
		fmt.Println(dbErr)
		fmt.Println("User not found for adding credentials.")
		return cred, dbErr
	}
	return cred, nil
}

func GetCredentialByUserId(userId uint64) ([]Credential, error) {
	scope := config.GetDefaultScope()
	var getCredential []Credential
	queryString := "SELECT x.credentials FROM `webauthn-bucket`.`webauthn-scope`.`users` x WHERE META().id=$userId"
	result, dbErr := config.ExecuteDBQuery(scope, queryString, &config.DBQueryParameters{
		"userId": userId,
	})
	if dbErr != nil {
		fmt.Println(dbErr)
		fmt.Println("User not found.")
		return getCredential, dbErr
	}

	err := result.One(&getCredential)
	if err != nil {
		fmt.Println(err)
		fmt.Println("User not found.")
		return getCredential, err
	}
	return getCredential, nil
}

func GetCredentialByCredId(credId []byte) (Credential, error) {
	var getCredential Credential
	result := db.Where("cred_id=?", credId).First(&getCredential)
	if result.Error != nil {
		return getCredential, result.Error
	}
	return getCredential, nil
}

// func UpdateCredentialByUserId(userId uint64, id, publicKey []byte, attestaionType string, authenticator Authenticator) (Credential, error) {
// 	var updateCredential Credential
// 	result := db.Model(&updateCredential).Where("user_id=?", userId).Updates(Credential{
// 		UserID:          userId,
// 		CredID:          id,
// 		PublicKey:       publicKey,
// 		AttestationType: attestaionType,
// 		Authenticator:   authenticator,
// 	})
// 	if result.Error != nil {
// 		return updateCredential, result.Error
// 	}
// 	return updateCredential, nil
// }
