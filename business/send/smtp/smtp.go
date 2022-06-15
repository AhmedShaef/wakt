// Package send provides a simple SMTP client.
package send

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"time"
)

// Email sends an email using the SMTP client based on the provided configuration.
func Email(from, to, subject, body string) error {

	server := mail.NewSMTPClient()
	server.Host = "smtp.mailtrap.io"
	server.Port = 587
	server.Username = "690fc903b3602b"
	server.Password = "648ad3fe72458f"
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).AddTo(to).SetSubject(subject).SetBody(mail.TextHTML, body)

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil

}
