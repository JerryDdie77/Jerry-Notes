package service_test

import (
	"context"
	"jerry-notes/config"
	"jerry-notes/internal/service"
	"testing"
)

func TestSendConfimationCode(t *testing.T) {

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	email := "shahterchannel@gmail.com"
	code := "123456"
	svc := service.NewEmailService(cfg.GmailToken)
	if err := svc.SendConfirmationCode(ctx, email, code); err != nil {
		t.Fatalf("SendConfirmationCode: %v", err)
	}

}
