package models

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/remaster/webauthn/pkg/config"

	"gorm.io/gorm"
)

var db *gorm.DB
var User webauthn.User

// type User struct {
// 	id []byte
// }

// Initialization function
func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(User)
}

// // User Model functions
// func (user *User) WebAuthnID() []byte {
// 	return user.id
// }

// func (user *User) WebAuthnName() string {
// 	return "newUser"
// }

// func (user *User) WebAuthnDisplayName() string {
// 	return "New User"
// }

// func (user *User) WebAuthnIcon() string {
// 	return "https://pics.com/avatar.png"
// }

// func (user *User) WebAuthnCredentials() []webauthn.Credential {
// 	return []webauthn.Credential{}
// }

// Database data extraction functions
func GetUsers() []webauthn.User {
	var users []webauthn.User
	db.Find(&users)
	return users
}

func CreateUser() *webauthn.User {
	db.Create(User)
	return &User
}

func GetUserById(Id int64) *webauthn.User {
	var getUser webauthn.User
	db.Where("ID=?", Id).Find(&getUser)
	return &getUser
}

func DeleteUser(Id int64) webauthn.User {
	deleteUser := GetUserById(Id)
	db.Delete(&deleteUser)
	return *deleteUser
}
