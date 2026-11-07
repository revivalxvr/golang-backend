package template

import (
	"email-service/src/common/dto/request"
	_ "embed"
	"html/template"
	"strings"

	"email-service/src/common/log"
)

//go:embed html/otp.html
var otpEmbed string

func NewOtp(data *request.SendOtp) *strings.Builder {
	t, err := template.New("otp").Parse(otpEmbed)
	if err != nil {
		log.Logger.Error(err.Error())
		return nil
	}

	var body strings.Builder

	if err := t.Execute(&body, data); err != nil {
		log.Logger.Error(err.Error())
		return nil
	}

	return &body
}
