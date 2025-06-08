package people

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func enrichPerson(person *Person) chan struct{} {
	promise := make(chan struct{})

	counter := make(chan struct{})
	go func() {
		person.Age = getAge(person.Name)
		counter <- struct{}{}
	}()
	go func() {
		person.Gender = getGender(person.Name)
		counter <- struct{}{}
	}()
	go func() {
		person.Nationality = getFullCountryName(person.Name)
		counter <- struct{}{}
	}()

	go func() {
		for range 3 {
			<-counter
		}
		close(counter)
		promise <- struct{}{}
		close(promise)
	}()

	return promise
}

func getAge(name string) uint8 {
	request := fmt.Sprintf("https://api.agify.io?name=%s", name)
	response, requestErr := http.Get(request)
	if requestErr != nil {
		log.Println(requestErr)
	}
	if response.StatusCode != 200 {
		log.Println(response.StatusCode)
		return 0
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
	if response.StatusCode != 200 {
		log.Println(response.StatusCode)
		return ""
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
	if response.StatusCode != http.StatusOK {
		log.Println(response.StatusCode)
		return ""
	}
	if requestErr != nil {
		log.Println(requestErr)
		return ""
	}
	var countries NationalityRequest
	decodeErr := json.NewDecoder(response.Body).Decode(&countries)
	if decodeErr != nil {
		log.Println(decodeErr)
		return ""
	}
	var mostProbablyCountry string
	var probability = .0

	for _, country := range countries.Countries {
		if country.Probability > probability {
			mostProbablyCountry = country.CountryID
			probability = country.Probability
		}
	}

	return mostProbablyCountry
}

func getFullCountryName(name string) string {
	code := getMostProbablyCountryCode(name)
	if code == "" {
		log.Println("country not found")
		return ""
	}
	request := fmt.Sprintf("https://restcountries.com/v3.1/alpha/%s", code)
	response, requestErr := http.Get(request)
	if requestErr != nil {
		log.Println(requestErr)
		return ""
	}
	var countriesInfo []CountryInfo
	decodeErr := json.NewDecoder(response.Body).Decode(&countriesInfo)
	if decodeErr != nil {
		log.Println(decodeErr)
	}

	return countriesInfo[0].Name.Common
}
