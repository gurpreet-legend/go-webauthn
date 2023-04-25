package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/remaster/webauthn/pkg/utils"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/remaster/webauthn/pkg/config"

	"gorm.io/gorm"
)

var db *gorm.DB

// var User User

type User struct {
	Id          uint64 `gorm:"primaryKey"`
	Name        string
	DisplayName string
	credentials []webauthn.Credential `gorm:"-:migration, type:text"`
}

// User Model functions
// WebAuthnID returns the user's ID
func (user User) WebAuthnID() []byte {
	id := utils.ConvertIntToByteArray(user.Id)
	return id
}

// WebAuthnName returns the user's username
func (user User) WebAuthnName() string {
	return user.Name
}

// WebAuthnDisplayName returns the user's display name
func (user User) WebAuthnDisplayName() string {
	return user.DisplayName
}

// WebAuthnIcon is not (yet) implemented
func (user User) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (user User) AddCredential(cred webauthn.Credential) {
	user.credentials = append(user.credentials, cred)
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
	rand.Seed(time.Now().UnixNano())
	user := &User{}
	user.Id = rand.Uint64()
	user.Name = name
	user.DisplayName = displayName
	fmt.Println(user.Id)
	fmt.Println(user.Name)
	fmt.Println(user.DisplayName)
	result := db.Create(&user)
	if result.Error != nil {
		fmt.Print(result.Error)
	}
	return *user
}

func GetUserByName(username string) (User, error) {
	var getUser User
	result := db.Where("name=?", username).First(&getUser)
	if result.Error != nil {
		return getUser, result.Error
	}
	return getUser, nil
}

func DeleteUser(username string) User {
	deleteUser, _ := GetUserByName(username)
	db.Delete(deleteUser)
	return deleteUser
}
