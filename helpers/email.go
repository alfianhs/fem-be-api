package helpers

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Mailer interface {
	From(fromEmail, fromName string)
	To(receiver []string)
	Subject(value string)
	Body(value string)
	Attachment(file []byte, filename string, contentType string)
	Send() error
}

type smtpMailer struct {
	email  *gomail.Message
	dialer *gomail.Dialer
}

func (mailer *smtpMailer) From(fromEmail, fromName string) {
	mailer.email.SetHeader("From", fmt.Sprintf("%s <%s>", fromName, fromEmail))
}

func (mailer *smtpMailer) To(val []string) {
	mailer.email.SetHeader("To", val...)
}

func (mailer *smtpMailer) Subject(val string) {
	mailer.email.SetHeader("Subject", val)
}

func (mailer *smtpMailer) Body(val string) {
	mailer.email.SetBody("text/html", val)
}

func (mailer *smtpMailer) Attachment(file []byte, filename string, c string) {
	mailer.email.Attach(
		filename,
		gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(file)
			return err
		}),
		gomail.SetHeader(map[string][]string{"Content-Type": {c}}),
	)
}

func (mailer *smtpMailer) Send() error {
	return mailer.dialer.DialAndSend(mailer.email)
}

func NewSMTPMailer() Mailer {
	// init mail
	mail := gomail.NewMessage()

	mailport, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		logrus.Error("Invalid Mail Port: ", err)
	}
	d := gomail.NewDialer(
		os.Getenv("MAIL_HOST"),
		mailport,
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_PASSWORD"),
	)

	mailer := smtpMailer{
		email:  mail,
		dialer: d,
	}

	// default sender
	senderEmail := os.Getenv("MAIL_FROM_ADDRESS")
	if senderEmail == "" {
		senderEmail = "admin@fem.id"
	}
	senderName := os.Getenv("MAIL_FROM_NAME")
	if senderName == "" {
		senderName = "Admin FEM"
	}
	mailer.From(senderEmail, senderName)

	return &mailer
}
