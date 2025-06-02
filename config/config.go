package config

import (
	"log"
	"os"
)

type Config struct {
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
}

func NewConfig() *Config {
	return &Config{
		InfoLogger:  log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),
		DebugLogger: log.New(os.Stdout, "[DEBUG]\t", log.Ldate|log.Ltime),
	}
}
