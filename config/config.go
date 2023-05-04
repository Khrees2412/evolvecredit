package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Println("No .env file found")
	}
	logrus.Println("Environment variables successfully loaded. Starting application...")
}
