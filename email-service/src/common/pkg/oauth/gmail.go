package oauth

import (
	"context"
	"email-service/env"
	"email-service/src/common/log"

	"google.golang.org/api/option"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func NewGmailService() *gmail.Service {
	config := &oauth2.Config{
		ClientID:     env.Conf.Oauth.GemailClientId,
		ClientSecret: env.Conf.Oauth.GmailClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{gmail.GmailSendScope},
	}
	token := &oauth2.Token{RefreshToken: env.Conf.Oauth.GmailRefreshToken}
	tokenSource := config.TokenSource(context.Background(), token)

	srvc, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	return srvc
}
