package controllers

import (
	"encoding/base32"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/remaster/webauthn/pkg/models"
	"github.com/remaster/webauthn/pkg/utils"
)

type TotpNonce struct {
	Nonce      string    `json:"nonce"`
	ExpiryTime time.Time `json:"expiryTime"`
}

func totpFallback() {
	fmt.Println("TOTP FALLBACK TRIGGERED")
}

func SaltAsB32(salt string) string {
	return base32.StdEncoding.EncodeToString(
		[]byte(salt),
	)
}

func VerifySaltFromOTP(passcode string, salt string, t time.Time) (bool, error) {

	salt32 := SaltAsB32(salt)

	return totp.ValidateCustom(
		passcode, salt32, t, totp.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    otp.DigitsEight,
			Algorithm: otp.AlgorithmSHA512,
		},
	)
}

func GenerateOTPFromSalt(salt string, t time.Time) (string, error) {

	slat32 := SaltAsB32(salt)
	return totp.GenerateCodeCustom(slat32, t, totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsEight,
		Algorithm: otp.AlgorithmSHA512,
	})
}

func GenerateOTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	user, err := models.GetUserByName(username)
	if err != nil {
		if err == models.ErrUserNotFound {
			errorMessage := fmt.Sprintf("User '%s' not found.", username)
			fmt.Println(errorMessage)
			utils.JsonResponse(w, errorMessage, http.StatusNotFound)
			return
		}

		// Other error occurred
		fmt.Println(err)
		utils.JsonResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	passcode, err := GenerateOTPFromSalt(strconv.FormatUint(user.Id, 10), time.Now())

	if err != nil {
		fmt.Println("Error generating totp")
		utils.JsonResponse(w, err, http.StatusInternalServerError)
		return
	}

	expiryTime := time.Now().Add(time.Second * time.Duration(30))
	res := TotpNonce{
		Nonce:      passcode,
		ExpiryTime: expiryTime,
	}
	fmt.Printf("Successfully generated totp `%s` with an expiry time of `%v`", passcode, expiryTime)
	// emailBody := fmt.Sprintf("OTP: %s and ExpiryTime: %s\n", passcode, expiryTime)
	// err = utils.SendEmailViaSMTP([]string{username}, []string{}, "Here is your webauthn OTP", emailBody, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	utils.JsonResponse(w, res, http.StatusOK)
}

type VerifyRequest struct {
	Nonce string `json:"nonce"`
}

func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	user, err := models.GetUserByName(username)
	if err != nil {
		if err == models.ErrUserNotFound {
			errorMessage := fmt.Sprintf("User '%s' not found.", username)
			fmt.Println(errorMessage)
			utils.JsonResponse(w, errorMessage, http.StatusNotFound)
			return
		}

		// Other error occurred
		fmt.Println(err)
		utils.JsonResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var responseBody VerifyRequest
	_ = json.NewDecoder(r.Body).Decode(&responseBody)
	verified, err := VerifySaltFromOTP(responseBody.Nonce, strconv.FormatUint(user.Id, 10), time.Now())
	if err != nil {
		fmt.Println("Error verifying totp")
		utils.JsonResponse(w, err, http.StatusInternalServerError)
		return
	}
	fmt.Printf("Successfully verified totp `%s`", responseBody.Nonce)
	utils.JsonResponse(w, verified, http.StatusOK)
}
