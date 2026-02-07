package service

import (
	"context"
	"log"
	"net/smtp"
)

type EmailService struct {
	gmailToken string
}

func NewEmailService(gmailToken string) *EmailService {
	return &EmailService{gmailToken: gmailToken}
}

func (e *EmailService) SendConfirmationCode(ctx context.Context, email, code string) error {

	auth := smtp.PlainAuth(
		"",
		"jerrynotesbusiness@gmail.com",
		e.gmailToken,
		"smtp.gmail.com",
	)

	msg := "Subject: Ваш код подтверждения\r\n" +
		"From: jerrynotesbusiness@gmail.com\r\n" +
		"\r\n" +
		"Ваш код подтверждения: " + code + "\r\n"

	err := smtp.SendMail("smtp.gmail.com:587",
		auth,
		"jerrynotesbusiness@gmail.com",
		[]string{email},
		[]byte(msg),
	)

	if err != nil {
		log.Printf("SendMail to %s: %v", email, err)
		return ErrInternal
	}

	return nil
}
