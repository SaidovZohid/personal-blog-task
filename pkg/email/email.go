package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/SaidovZohid/personal-blog-task/config"
)

type SendEmailRequest struct {
	To      []string
	Type    string
	Body    map[string]string
	Subject string
}

const (
	VerificationEmail = "verification_email"
)

func SendEmail(cfg *config.Config, req *SendEmailRequest) error {
	fmt.Println("I am here 1")
	from := cfg.Smtp.Sender
	to := req.To

	password := cfg.Smtp.Password

	var body bytes.Buffer

	templatePath := getTemplatePath(req.Type)
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	fmt.Println("I am here 2")
	t.Execute(&body, req.Body)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := fmt.Sprintf("Subject: %s\n", req.Subject)
	msg := []byte(subject + mime + body.String())

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	fmt.Println("I am here 3")
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	fmt.Println("I am here 4", err)
	if err != nil {
		return err
	}
	return nil
}

func getTemplatePath(emailType string) string {
	switch emailType {
	case VerificationEmail:
		return "./templates/verification_email.html"
	}
	return ""
}
