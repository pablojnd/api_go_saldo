package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv carga las variables de entorno desde el archivo .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error al cargar el archivo .env: %v", err)
	}
}

// GetEnv obtiene el valor de una variable de entorno o devuelve un valor por defecto
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
