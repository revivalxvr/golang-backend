package service

import (
	"email-service/src/common/dto/request"
	v "email-service/src/common/pkg/validator"
	"email-service/src/factory"
	"email-service/template"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/gmail/v1"
)

type RabbitMQ struct {
	gmailService *gmail.Service
}

func NewRabbitMQ(f *factory.Factory) *RabbitMQ {
	return &RabbitMQ{
		gmailService: f.GmailService,
	}
}

func (s *RabbitMQ) SendOtp(data []byte) error {
	req := new(request.SendOtp)

	if err := json.Unmarshal(data, req); err != nil {
		return err
	}
	req.Year = time.Now().Year()

	fmt.Printf("OTP Request: %+v\n", req)

	if err := v.Validate.Struct(req); err != nil {
		return err
	}
	tmpl := template.NewOtp(req)

	emailTo := "To :" + req.Email + "\r\n"
	emailSubject := "Subject: OTP Verification\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + emailSubject + mime + "\n" + tmpl.String())

	m := new(gmail.Message)
	m.Raw = base64.StdEncoding.EncodeToString(msg)

	_, err := s.gmailService.Users.Messages.Send("me", m).Do()
	if err != nil {
		return err
	}
	return nil
}
