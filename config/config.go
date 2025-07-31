package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port               string
	DatabaseURL        string
	GoogleClientID     string
	GoogleClientSecret string
	MobileAppSchema    string
	BaseURL            string
	AppEnv             string // development, qa, or production
	JWTSecret          string
	JWTRefreshSecret   string
}

func SetupAppConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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

	if googleClientId == "" || googleClientSecret == "" {
		log.Fatal("env vars GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET are not set")
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

	return &AppConfig{
		Port:               port,
		DatabaseURL:        dbUrl,
		GoogleClientID:     googleClientId,
		GoogleClientSecret: googleClientSecret,
		MobileAppSchema:    mobileAppSchema,
		BaseURL:            baseUrl,
		AppEnv:             appEnv,
	}
}
