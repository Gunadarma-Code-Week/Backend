package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"
)

type EmailService struct {
	auth smtp.Auth
	from string
	addr string
}

func NewEmailService() *EmailService {
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	from := os.Getenv("EMAIL_FROM")
	port := os.Getenv("EMAIL_PORT")

	return &EmailService{
		auth: smtp.PlainAuth("", username, password, host),
		from: from,
		addr: fmt.Sprintf("%s:%s", host, port),
	}
}

func (mh *EmailService) SendEmail(
	subject string,
	to []string,
	msg string,
) error {
	body := "Subject: " + subject + "!\n"
	body += "From: " + mh.from + "\n"
	body += "To: " + strings.Join(to, ",") + "\n"
	body += "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	body += msg

	if err := smtp.SendMail(mh.addr, mh.auth, mh.from, to, []byte(body)); err != nil {
		return err
	}
	return nil
}

func (mh *EmailService) SendEmailHTML(
	subject string,
	to []string,
	templateFileName string,
	data any,
) error {
	body := "Subject: " + subject + "!\n"
	body += "From: " + mh.from + "\n"
	body += "To: " + strings.Join(to, ",") + "\n"
	body += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	body += buf.String()

	if err := smtp.SendMail(mh.addr, mh.auth, mh.from, to, []byte(body)); err != nil {
		return err
	}
	return nil
}
