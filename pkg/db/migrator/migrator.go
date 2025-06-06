package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	envLoadErr := godotenv.Load("../../../.env")
	if envLoadErr != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	ssl := os.Getenv("DB_SSL")

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user, password, host, port, dbname, ssl,
	)

	db, openErr := sql.Open("postgres", dsn)
	if openErr != nil {
		log.Fatal("DB connection error", openErr)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Fatal(closeErr)
		}
	}()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS people (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            surname TEXT NOT NULL,
            patronymic TEXT,
            age INTEGER NOT NULL CHECK (age >= 0),
            gender TEXT NOT NULL,
            nationality TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW() 
        )`,
	}

	for _, q := range queries {
		fmt.Println("in progress...:\n", q)
		_, err := db.Exec(q)
		if err != nil {
			log.Fatal("migration error", err)
		}
	}

	fmt.Println("All migrations have been applied.")
}
