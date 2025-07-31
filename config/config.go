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

	return &AppConfig{
		Port:               port,
		DatabaseURL:        dbUrl,
		GoogleClientID:     googleClientId,
		GoogleClientSecret: googleClientSecret,
	}
}
