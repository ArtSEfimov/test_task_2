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

type Country struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
type NationalityRequest struct {
	Countries []Country `json:"countries"`
}

type DetailInfo struct {
	Common string `json:"common"`
}
type CountryInfo struct {
	Name DetailInfo `json:"name"`
}
