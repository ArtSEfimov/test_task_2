package config

import (
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
	return &Config{
		InfoLogger:  log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),
		DebugLogger: log.New(os.Stdout, "[DEBUG]\t", log.Ldate|log.Ltime),
	}
}
