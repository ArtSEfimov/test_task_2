package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type DB struct {
	Driver   string
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	SSL      string
}

type Config struct {
	Port        string
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
	Database    DB
}

func NewConfig() *Config {
	envLoadErr := godotenv.Load("../../.env")
	if envLoadErr != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	// DB
	driver := os.Getenv("DB_DRIVER")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	ssl := os.Getenv("DB_SSL")

	return &Config{
		Port:        port,
		InfoLogger:  log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),
		DebugLogger: log.New(os.Stdout, "[DEBUG]\t", log.Ldate|log.Ltime),
		Database: DB{
			Driver:   driver,
			User:     user,
			Password: password,
			Host:     host,
			Port:     dbPort,
			Name:     dbName,
			SSL:      ssl,
		},
	}
}
