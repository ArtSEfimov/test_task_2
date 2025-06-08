package people

import (
	"database/sql"
	"errors"
	"fmt"
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

func (repository *Repository) Get(query string, people *AllPeopleResponse, params ...any) error {
	rows, queryErr := repository.Database.DB.Query(query, params...)
	if queryErr != nil {
		return queryErr
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	for rows.Next() {
		var person Person
		scanErr := rows.Scan(
			&person.ID,
			&person.CreatedAt,
			&person.UpdatedAt,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality,
		)

		if scanErr != nil {
			return scanErr
		}

		people.People = append(people.People, person)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return rowsErr
	}

	return nil
}

func (repository *Repository) GetByID(query string, person *Person, id uint64) error {
	row := repository.Database.DB.QueryRow(query, id)

	scanErr := row.Scan(
		&person.ID,
		&person.CreatedAt,
		&person.UpdatedAt,
		&person.Name,
		&person.Surname,
		&person.Patronymic,
		&person.Age,
		&person.Gender,
		&person.Nationality,
	)

	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return fmt.Errorf("record with id %d not found", id)
		}
		return scanErr
	}
	return nil
}

func (repository *Repository) Create(query string, person *Person) error {
	queryErr := repository.Database.DB.QueryRow(
		query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
	).Scan(&person.ID, &person.CreatedAt, &person.UpdatedAt)

	if queryErr != nil {
		return queryErr
	}
	return nil
}

func (repository *Repository) Update(query string, person *Person, id uint64) error {
	queryErr := repository.Database.DB.QueryRow(
		query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
		id,
	).Scan(&person.ID, &person.CreatedAt, &person.UpdatedAt)

	if queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			return fmt.Errorf("record with id %d not found", id)
		}
		return queryErr
	}
	return nil
}

func (repository *Repository) Delete(query string, id uint64) error {
	result, queryErr := repository.Database.DB.Exec(query, id)
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("record with id %d not found", id)
	}

	if queryErr != nil {
		return queryErr
	}
	return nil
}
