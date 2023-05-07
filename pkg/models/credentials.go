package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/remaster/webauthn/pkg/config"
)

type Authenticator struct {
	AAGUID       []byte `gorm:"aaguid"` // []byte
	SignCount    uint32
	CloneWarning bool
}

type Credential struct {
	UserID          uint64
	CredID          []byte `gorm:"size:700;primaryKey"` // []byte
	PublicKey       []byte // []byte
	AttestationType string
	Authenticator   Authenticator `gorm:"foreignKey:AAGUID"`
}

// Initialization function
func init() {
	config.ConnectDB()
	db = config.GetDB()
	db.AutoMigrate(&Authenticator{})
	db.AutoMigrate(&Credential{})
}

func CreateCredential(userId uint64, id, publicKey []byte, attestaionType string, authenticator Authenticator) Credential {
	rand.Seed(time.Now().UnixNano())
	cred := &Credential{}
	cred.UserID = userId
	cred.CredID = id
	cred.PublicKey = publicKey
	cred.AttestationType = attestaionType
	cred.Authenticator = authenticator
	result := db.Create(&cred)
	if result.Error != nil {
		fmt.Print(result.Error)
	}
	return *cred
}

func GetCredentialByUserId(userId uint64) ([]Credential, error) {
	var getCredential []Credential
	result := db.Where("user_id=?", userId).Preload("Authenticator").Find(&getCredential)
	if result.Error != nil {
		return getCredential, result.Error
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

func UpdateCredentialByUserId(userId uint64, id, publicKey []byte, attestaionType string, authenticator Authenticator) (Credential, error) {
	var updateCredential Credential
	result := db.Model(&updateCredential).Where("user_id=?", userId).Updates(Credential{
		UserID:          userId,
		CredID:          id,
		PublicKey:       publicKey,
		AttestationType: attestaionType,
		Authenticator:   authenticator,
	})
	if result.Error != nil {
		return updateCredential, result.Error
	}
	return updateCredential, nil
}
