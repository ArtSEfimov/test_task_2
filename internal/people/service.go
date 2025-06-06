package people

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func enrichPerson(person *Person) {
	person.Age = getAge(person.Name)
	person.Gender = getGender(person.Name)
	person.Nationality = getFullCountryName(person.Name)
}

func getAge(name string) uint8 {
	request := fmt.Sprintf("https://api.agify.io?name=%s", name)
	response, requestErr := http.Get(request)
	if requestErr != nil {
		log.Println(requestErr)
	}
	var age AgeRequest
	decodeErr := json.NewDecoder(response.Body).Decode(&age)
	if decodeErr != nil {
		log.Println(decodeErr)
	}
	return age.Age
}

func getGender(name string) string {
	request := fmt.Sprintf("https://api.genderize.io?name=%s", name)
	response, requestErr := http.Get(request)
	if requestErr != nil {
		log.Println(requestErr)
	}
	var gender GenderRequest
	decodeErr := json.NewDecoder(response.Body).Decode(&gender)
	if decodeErr != nil {
		log.Println(decodeErr)
	}
	return gender.Gender
}

func getMostProbablyCountryCode(name string) string {
	request := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	response, requestErr := http.Get(request)
	if requestErr != nil {
		log.Println(requestErr)
	}
	var countries NationalityRequest
	decodeErr := json.NewDecoder(response.Body).Decode(&countries)
	if decodeErr != nil {
		log.Println(decodeErr)
	}
	var mostProbablyCountry string
	var probability = .0
	for _, country := range countries.Countries {
		if country.Probability > probability {
			mostProbablyCountry = country.CountryID
		}
	}

	return mostProbablyCountry
}

func getFullCountryName(name string) string {
	code := getMostProbablyCountryCode(name)
	request := fmt.Sprintf("https://restcountries.com/v3.1/alpha/%s", code)
	response, requestErr := http.Get(request)
	if requestErr != nil {
		log.Println(requestErr)
	}
	var countriesInfo []CountryInfo
	decodeErr := json.NewDecoder(response.Body).Decode(&countriesInfo)
	if decodeErr != nil {
		log.Println(decodeErr)
	}

	return countriesInfo[0].Name.Common
}
func validateData()
