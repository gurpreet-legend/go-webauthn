package utils

import (
	"crypto/tls"
	"log"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func SendEmailViaSMTP(tos []string, cc []string, title string, body string, processFunc func(string, string, int) string) error {

	client := mail.NewSMTPClient()
	client.Host = "smtp.gmail.com"
	client.Port = 587
	client.Username = "lostymailtest@gmail.com"
	client.Password = "Testadmin123"
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
		email = email.SetFrom("lostymailtest@gmail.com").SetSubject(title).AddTo(to)
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
