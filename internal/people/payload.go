package people

type Request struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type AllPeopleResponse struct {
	People []Person `json:"people"`
}

type AgeRequest struct {
	Age uint8 `json:"age"`
}
type GenderRequest struct {
	Gender string `json:"gender"`
}
