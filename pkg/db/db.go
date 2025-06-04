package db

import (
	"database/sql"
	"fmt"
	"go_test_task_2/config"
)

type DB struct {
	*sql.DB
}

func NewDB(config *config.Config) *DB {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name, config.Database.SSL)
	driver := config.Database.Driver
	db, openErr := sql.Open(driver, dsn)
	if openErr != nil {
		panic(openErr)
	}
	return &DB{db}
}
