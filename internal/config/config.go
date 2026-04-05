package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DSN  string
}

func Load() *Config {
	var err = godotenv.Load()
	var user = os.Getenv("DB_USER")
	var password = os.Getenv("DB_PASSWORD")
	var dbName = os.Getenv("DB_NAME")
	var dbHost = os.Getenv("DB_HOST")
	var dbPort = os.Getenv("DB_PORT")
	var dbSslmode = os.Getenv("DB_SSLMODE")
	if err != nil {
		log.Println("Error loading .env file")
	}

	dsn := fmt.Sprintf(`user=%s password=%s dbname=%s host=%s port=%s sslmode=%s`, user, password, dbName, dbHost, dbPort, dbSslmode)

	return &Config{
		Port: os.Getenv("APP_PORT"),
		DSN:  dsn,
	}
}
