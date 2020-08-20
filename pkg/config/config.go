package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	BotOAuthToken     string
	VerificationToken string
	SqliteFilePath    string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		BotOAuthToken:     os.Getenv("SLACK_BOT_OAUTH_TOKEN"),
		VerificationToken: os.Getenv("SLACK_VERIFICATION_TOKEN"),
		SqliteFilePath:    os.Getenv("SQLITE_FILE_PATH"),
	}
}
