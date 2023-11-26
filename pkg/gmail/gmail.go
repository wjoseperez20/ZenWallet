package gmail

import (
	"gopkg.in/gomail.v2"
	"os"
)

var Mailer *gomail.Dialer

func ConnectGmail() {

	gmailUser := os.Getenv("GMAIL_USER")
	gmailSecret := os.Getenv("GMAIL_SECRET")

	Mailer = gomail.NewDialer("smtp.gmail.com", 587, gmailUser, gmailSecret)
	Mailer.TLSConfig = nil
}
