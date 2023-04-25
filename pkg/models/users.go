package models

import (
	"fmt"

	"github.com/remaster/webauthn/pkg/utils"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/remaster/webauthn/pkg/config"

	"gorm.io/gorm"
)

var db *gorm.DB

// var User User

type User struct {
	gorm.Model
	id          uint64
	name        string
	displayName string
	credentials []webauthn.Credential
}

// User Model functions
// WebAuthnID returns the user's ID
func (user User) WebAuthnID() []byte {
	id := utils.ConvertIntToByteArray(user.id)
	return id
}

// WebAuthnName returns the user's username
func (user User) WebAuthnName() string {
	return user.name
}

// WebAuthnDisplayName returns the user's display name
func (user User) WebAuthnDisplayName() string {
	return user.displayName
}

// WebAuthnIcon is not (yet) implemented
func (user User) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u User) AddCredential(cred webauthn.Credential) {
	u.credentials = append(u.credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (user User) WebAuthnCredentials() []webauthn.Credential {
	return user.credentials
}

// Initialization function
func init() {
	config.ConnectDB()
	db = config.GetDB()
	db.AutoMigrate(&User{})
}

// Database data extraction functions
func GetUsers() []User {
	var users []User
	db.Find(&users)
	return users
}

func CreateUser(name string, displayName string) User {
	user := &User{}
	user.id = utils.RandomUint64()
	user.name = name
	user.displayName = displayName
	result := db.Create(&user)
	if result.Error != nil {
		fmt.Print(result.Error)
	}
	return *user
}

func GetUserById(Id int64) *User {
	var getUser User
	db.Where("ID=?", Id).Find(&getUser)
	return &getUser
}

func DeleteUser(Id int64) User {
	deleteUser := GetUserById(Id)
	db.Delete(&deleteUser)
	return *deleteUser
}
