package models

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/remaster/webauthn/pkg/utils"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/remaster/webauthn/pkg/config"

	"gorm.io/gorm"
)

var db *gorm.DB

// var User User

type User struct {
	Id          uint64 `gorm:"primaryKey"`
	Name        string
	DisplayName string
	Credentials []webauthn.Credential `gorm:"-"`
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
func (user *User) AddCredential(cred webauthn.Credential) {
	user.Credentials = append(user.Credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (user User) WebAuthnCredentials() []webauthn.Credential {
	return user.Credentials
}

// Initialization function
func init() {
	config.ConnectDB()
	db = config.GetDB()
	db.AutoMigrate(&User{})
}

type CouchUser struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// Couchbase functions
func CouchGetUsers() ([]CouchUser, error) {
	scope := config.GetDefaultScope()
	queryString := "SELECT x.* FROM `webauthn-bucket`.`webauthn-scope`.`users` x"
	var users []CouchUser
	result, err := config.ExecuteDBQuery(scope, queryString, &config.DBQueryParameters{})
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error getting all the users.")
		return users, err
	}

	for result.Next() {
		var user CouchUser
		err := result.Row(&user)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error getting all the users.")
			return users, err
		}

		users = append(users, user)
	}
	return users, nil
}

func CouchGetUserByName(username string) (CouchUser, error) {
	scope := config.GetDefaultScope()
	var user CouchUser
	queryString := "SELECT x.* FROM `webauthn-bucket`.`webauthn-scope`.`users` x WHERE x.name=$username"
	result, dbErr := config.ExecuteDBQuery(scope, queryString, &config.DBQueryParameters{
		"username": username,
	})
	if dbErr != nil {
		fmt.Println(dbErr)
		fmt.Println("Error getting all the users.")
		return user, dbErr
	}

	err := result.One(&user)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error getting all the users.")
		return user, err
	}
	return user, nil
}

func CouchCreateUser(name string, displayName string) error {
	rand.Seed(time.Now().UnixNano())
	user := &CouchUser{}
	user.Id = rand.Uint64()
	user.Name = name
	user.DisplayName = displayName

	users := config.GetDefaultScope().Collection("users")
	_, err := users.Insert(strconv.FormatUint(user.Id, 10), user, &gocb.InsertOptions{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error creating user.")
		return err
	}
	return nil
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
