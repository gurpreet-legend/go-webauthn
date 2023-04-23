package models

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/remaster/webauthn/pkg/config"

	"gorm.io/gorm"
)

var db *gorm.DB
var User webauthn.User

// type User struct {

// }

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(User)
}

func GetUsers() []webauthn.User {
	var users []webauthn.User
	db.Find(&users)
	return users
}

func (b *webauthn.User) CreateUser() *webauthn.User {
	db.Create(&b)
	return b
}

func GetUserById(Id int64) (*webauthn.User, *gorm.DB) {
	var getUser webauthn.User
	db := db.Where("ID=?", Id).Find(&getUser)
	return &getUser, db
}

func DeleteUser(Id int64) webauthn.User {
	deleteUser, _ := GetUserById(Id)
	db.Delete(&deleteUser)
	return *deleteUser
}
