package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/cirrostratus-cloud/common/email"
)

type SESEmailService struct {
	sesClient *ses.Client
}

func NewSESMailService(sesClient *ses.Client) email.EmailService {
	return &SESEmailService{sesClient: sesClient}
}

func (e *SESEmailService) SendEmail(from string, to string, subject string, body string) error {
	_, err := e.sesClient.SendEmail(
		context.TODO(),
		&ses.SendEmailInput{
			Destination: &types.Destination{
				ToAddresses: []string{to},
			},
			Message: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(body),
					},
				},
				Subject: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(subject),
				},
			},
			Source: aws.String(from),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
