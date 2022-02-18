package common

import (
	"fmt"
	"net/smtp"
)

const (
	fromAddress       = "pranotobudi.app@gmail.com"
	fromEmailPassword = "CamelCasePasswordBismillah"
	smtpServer        = "smtp.gmail.com"
	smtpPort          = "587"
)

// SendEmail will send email to toAddress with body as its content.
func SendEmail(toAddress []string, body string) {
	// Message.
	//   message := []byte("This is a test email message.")
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "Chat Messages Archive" + "!\n"
	message := []byte(subject + mime + "\n" + body)

	// Authentication.
	auth := smtp.PlainAuth("", fromAddress, fromEmailPassword, smtpServer)

	// Sending email.
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, fromAddress, toAddress, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}
