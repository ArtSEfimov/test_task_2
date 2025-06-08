package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const table = "people"

func main() {
	envLoadErr := godotenv.Load("../../.env")
	if envLoadErr != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	ssl := os.Getenv("DB_SSL")
	driver := os.Getenv("DB_DRIVER")

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user, password, host, port, dbName, ssl,
	)

	db, openErr := sql.Open(driver, dsn)
	if openErr != nil {
		log.Fatal("DB connection error", openErr)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Fatal(closeErr)
		}
	}()

	queries := []string{
		fmt.Sprintf(`DROP TABLE IF EXISTS %s`, table),
		fmt.Sprintf(`CREATE TABLE %s (
            id SERIAL PRIMARY KEY,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
            name TEXT NOT NULL,
            surname TEXT NOT NULL,
            patronymic TEXT,
            age INTEGER NOT NULL CHECK (age >= 0),
            gender TEXT NOT NULL,
            nationality TEXT NOT NULL
        )`, table),
		`CREATE OR REPLACE FUNCTION update_updated_at()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;`,
		fmt.Sprintf(`
    CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON %s
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();
    `, table),
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
