package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	User          string
	Password      string
	AdminUser     string
	AdminPass     string
	Admin2        string
	Pass2         string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := Config{
		TelegramToken: os.Getenv("BOT_TOKEN"),
		User:          os.Getenv("USER"),
		Password:      os.Getenv("PASSWORD"),
		AdminUser:     os.Getenv("ADMIN_USER"),
		AdminPass:     os.Getenv("ADMIN_PASS"),
		Admin2:        os.Getenv("ADMIN2"),
		Pass2:         os.Getenv("PASS2"),
	}

	return config
}
