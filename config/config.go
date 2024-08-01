package config

import (
	"log"
	"os"
)

var (
	BotToken  string
	Password  string
	User      string
	AdminUser string
	AdminPass string
	Admin2    string
	Pass2     string
)

func LoadConfig() {
	BotToken = getEnv("BOT_TOKEN", "default_bot_token")
	Password = getEnv("PASSWORD", "default_password")
	User = getEnv("USER", "default_user")
	AdminUser = getEnv("AdminUser", "default_adminuser")
	AdminPass = getEnv("AdminPass", "default_adminpass")
	Admin2 = getEnv("Admin2", "default_admin2")
	Pass2 = getEnv("Pass2", "default_pass2")

}

func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Environment variable %s not set, using default value: %s", key, fallback)
		return fallback
	}
	return value
}
