package people

type Person struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         uint8  `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}
