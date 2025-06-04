package people

import "time"

type DB struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Person struct {
	DB
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         uint8  `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}
