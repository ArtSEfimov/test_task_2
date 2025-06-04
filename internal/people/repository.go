package people

import (
	"go_test_task_2/pkg/db"
)

type Repository struct {
	Database *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{
		Database: database,
	}
}
