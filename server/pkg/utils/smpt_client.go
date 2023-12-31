package utils

import (
	"crypto/tls"
	"log"
	"os"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func SendEmailViaSMTP(tos []string, cc []string, title string, body string, processFunc func(string, string, int) string) error {

	client := mail.NewSMTPClient()
	client.Host = "smtp.gmail.com"
	client.Port = 587
	client.Username = os.Getenv("INFO_EMAIL")
	client.Password = os.Getenv("INFO_PASSWORD")
	client.Encryption = mail.EncryptionSTARTTLS
	client.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	client.ConnectTimeout = 60 * time.Second
	client.SendTimeout = 60 * time.Second
	client.KeepAlive = true

	sender, err := client.Connect()
	if err != nil {
		return err
	}

	for idx, to := range tos {
		email := mail.NewMSG()
		email = email.SetFrom("info@remaster.io").SetSubject(title).AddTo(to)
		email = email.AddCc(cc...)

		if processFunc != nil {
			bodyProcessed := processFunc(body, to, idx)
			email = email.SetBody(mail.TextHTML, bodyProcessed)
		} else {
			email = email.SetBody(mail.TextHTML, body)
		}

		if email.Error != nil {
			return err
		}

		err := email.Send(sender)

		if err != nil {
			return err
		}

		log.Printf("sent mail to %s", to)
	}

	return nil
}
