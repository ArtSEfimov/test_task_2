package people

import "time"

type DB struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Person struct {
	DB
	Name        string `json:"name" validate:"required"`
	Surname     string `json:"surname" validate:"required"`
	Patronymic  string `json:"patronymic"`
	Age         uint8  `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}
