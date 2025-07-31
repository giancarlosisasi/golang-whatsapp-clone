package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang-whatsapp-clone/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthService struct {
	config     *oauth2.Config
	appConfig  *config.AppConfig
	jwtService *JWTService
}

func NewOAuthService(appConfig *config.AppConfig, jwtService *JWTService) *OAuthService {
	googleOAuth := &oauth2.Config{
		ClientID:     appConfig.GoogleClientID,
		ClientSecret: appConfig.GoogleClientSecret,
		RedirectURL:  appConfig.BaseURL + "/api/v1/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &OAuthService{
		config:     googleOAuth,
		jwtService: jwtService,
		appConfig:  appConfig,
	}
}

func (o *OAuthService) GenerateAuthURL() (url string, state string, err error) {
	state, err = o.generateState()
	if err != nil {
		return "", "", err
	}

	url = o.config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	return url, state, nil
}

func (o *OAuthService) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return o.config.Exchange(ctx, code)
}

func (o *OAuthService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	client := o.config.Client(ctx, token)
	resp, err := client.Get(o.appConfig.GoogleUserInfoUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var googleUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &googleUser, nil
}

func (o *OAuthService) generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
