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
		log.Printf("объект с ID %d не обнаружен, %v", id, scanErr)
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
		return queryErr
	}
	return nil
}
func (repository *Repository) Delete(query string, id uint64) error {
	_, queryErr := repository.Database.DB.Exec(query, id)

	if queryErr != nil {
		log.Println("delete error: ", queryErr)
		return queryErr
	}
	return nil
}
