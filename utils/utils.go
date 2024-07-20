package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadToken() (string, string, string) {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Obtener el valor de una variable de entorno
	botToken := os.Getenv("BOT_TOKEN")
	password := os.Getenv("PASSWORD")
	user := os.Getenv("USER")

	return botToken, password, user
}
