package models

import (
	"errors"
	"fmt"
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
var ErrUserNotFound = errors.New("User not found")

type User struct {
	Id          uint64                `json:"id"`
	Name        string                `json:"name"`
	DisplayName string                `json:"displayName"`
	Credentials []webauthn.Credential `json:"credentials"`
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

// Couchbase functions
func GetUsers() ([]User, error) {
	scope := config.GetDefaultScope()
	queryString := "SELECT x.* FROM `webauthn-bucket`.`webauthn-scope`.`users` x"
	var users []User
	result, err := config.ExecuteDBQuery(scope, queryString, &config.DBQueryParameters{})
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error getting all the users.")
		return users, err
	}

	for result.Next() {
		var user User
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

func GetUserByName(username string) (User, error) {
	scope := config.GetDefaultScope()
	var user User
	queryString := "SELECT x.* FROM `webauthn-bucket`.`webauthn-scope`.`users` x WHERE x.name = $username"
	result, dbErr := config.ExecuteDBQuery(scope, queryString, &config.DBQueryParameters{
		"username": username,
	})
	if dbErr != nil {
		fmt.Println(dbErr)
		fmt.Println("User not found.")
		return user, dbErr
	}

	err := result.One(&user)
	if err != nil {
		fmt.Println(err)
		fmt.Println("User not found.")
		return user, ErrUserNotFound
	}
	return user, nil
}

func CreateUser(name string, displayName string) error {
	// rand.Seed(time.Now().UnixNano())
	// userId := rand.Uint64()
	userId := utils.GenerateUserID()
	user := &User{
		Id:          userId,
		Name:        name,
		DisplayName: displayName,
		Credentials: []webauthn.Credential{},
	}

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
