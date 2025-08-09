package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port               string
	DatabaseURL        string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleUserInfoUrl  string
	MobileAppSchema    string
	BaseURL            string
	AppEnv             string // development, qa, or production
	JWTSecret          string
	JWTRefreshSecret   string
	CookieName         string
}

func SetupAppConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		// TODO: only log this when we are not in production
		fmt.Printf("warning: error to load env: %s", err.Error())
	}

	viper.AutomaticEnv()

	dbUrl := viper.GetString("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("env var DATABASE_URL is not set")
	}
	port := viper.GetString("PORT")

	// google auth
	googleClientId := viper.GetString("GOOGLE_CLIENT_ID")
	googleClientSecret := viper.GetString("GOOGLE_CLIENT_SECRET")
	googleUserInfo := viper.GetString("GOOGLE_USER_INFO_URL")

	if googleClientId == "" || googleClientSecret == "" || googleUserInfo == "" {
		log.Fatal("env vars GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET or GOOGLE_USER_INFO_URL are not set")
	}

	mobileAppSchema := viper.GetString("MOBILE_APP_SCHEME")
	if mobileAppSchema == "" {
		log.Fatal("env vars MOBILE_APP_SCHEME is not set")
	}

	baseUrl := viper.GetString("BASE_URL")
	if baseUrl == "" {
		log.Fatal("env vars BASE_URL is not set")
	}

	appEnv := viper.GetString("APP_ENV")
	if appEnv == "" {
		log.Fatal("env vars APP_ENV is not set")
	}

	jwtSecret := viper.GetString("JWT_SECRET")
	jwtRefreshSecret := viper.GetString("JWT_REFRESH_SECRET")
	if jwtSecret == "" || jwtRefreshSecret == "" {
		log.Fatal("env vars JWT_SECRET or JWT_REFRESH_SECRET are not set")
	}

	cookieName := viper.GetString("COOKIE_NAME")
	if cookieName == "" {
		log.Fatal("env vars COOKIE_NAME is not set")
	}

	return &AppConfig{
		Port:               port,
		DatabaseURL:        dbUrl,
		GoogleClientID:     googleClientId,
		GoogleClientSecret: googleClientSecret,
		GoogleUserInfoUrl:  googleUserInfo,
		MobileAppSchema:    mobileAppSchema,
		BaseURL:            baseUrl,
		AppEnv:             appEnv,
		CookieName:         cookieName,
	}
}
