package people

import (
	"go_test_task_2/pkg/db"
	"log"
)

type Repository struct {
	Database *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{
		Database: database,
	}
}

func (repository *Repository) GetAll(query string, people *AllPeopleResponse) error {
	rows, queryErr := repository.Database.DB.Query(query)
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
		scanErr := rows.Scan(&person.ID,
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
			log.Fatal("Ошибка при сканировании строки: ", scanErr)
			return scanErr
		}

		people.People = append(people.People, person)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		log.Fatal("Ошибка при обработке строк: ", rowsErr)
		return rowsErr
	}

	return nil
}
