package config

import (
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	w   *webauthn.WebAuthn
	err error
)

func SetupWebAuthn() {
	//Webauthn setup
	wconfig := &webauthn.Config{
		RPDisplayName: "Go Webauthn",                     // Display Name for your site
		RPID:          "localhost",                       // Generally the FQDN for your site
		RPOrigins:     []string{"http://localhost:3000"}, // The origin URLs allowed for WebAuthn requests
	}

	if w, err = webauthn.New(wconfig); err != nil {
		fmt.Println(err)
	}
}

func GetWebAuthn() *webauthn.WebAuthn {
	return w
}
