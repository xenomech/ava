package email

import (
	"context"

	"ava/config"

	"github.com/resend/resend-go/v2"
)

type Service interface {
	Send(ctx context.Context, to, subject, templateName string, data any) error
}

type emailService struct {
	client    *resend.Client
	fromEmail string
	fromName  string
}

func NewService() Service {
	cfg := config.GetConfig()

	return &emailService{
		client:    resend.NewClient(cfg.ResendAPIKey),
		fromEmail: cfg.ResendFromEmail,
		fromName:  cfg.ResendFromName,
	}
}
