package utils

import (
	"crypto/tls"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

func SendEmail(receipantEmail string) error {
	emailMessage := gomail.NewMessage()
	emailMessage.SetHeader("From", os.Getenv("SENDING_EMAIL"))
	emailMessage.SetHeader("To", receipantEmail)
	emailMessage.SetBody("text/plain", "You have been assigned some tasks. Please register to access them.")
	portNumber, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	dialer := gomail.NewDialer("smtp.gmail.com", portNumber, os.Getenv("SENDING_EMAIL"), os.Getenv("SENDING_PASSWORD"))
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(emailMessage)
	return err
}
