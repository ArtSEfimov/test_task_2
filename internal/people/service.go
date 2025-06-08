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
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()

	if requestErr != nil {
		log.Printf("Error sending request to Agify API for name: %s, error: %v", name, requestErr)
		return 0
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("Received non-success status code from Agify API for name: %s, status code: %d", name, response.StatusCode)
		return 0
	}

	var age AgeRequest
	decodeErr := json.NewDecoder(response.Body).Decode(&age)
	if decodeErr != nil {
		log.Printf("Error decoding JSON response from Agify API for name: %s, error: %v", name, decodeErr)
		return 0
	}
	return age.Age
}

func getGender(name string) string {
	request := fmt.Sprintf("https://api.genderize.io?name=%s", name)
	response, requestErr := http.Get(request)
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()
	if requestErr != nil {
		log.Printf("Error sending request to Genderize API for name: %s, error: %v", name, requestErr)
		return ""
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("Received non-success status code from Genderize API for name: %s, status code: %d", name, response.StatusCode)
		return ""
	}
	var gender GenderRequest
	decodeErr := json.NewDecoder(response.Body).Decode(&gender)
	if decodeErr != nil {
		log.Printf("Error decoding JSON response from Genderize API for name: %s, error: %v", name, decodeErr)
		return ""
	}
	return gender.Gender
}

func getMostProbablyCountryCode(name string) string {
	request := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	response, requestErr := http.Get(request)
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()
	if response.StatusCode != http.StatusOK {
		log.Printf("Error: Received non-OK status code from Nationalize API for name: %s, status code: %d", name, response.StatusCode)
		return ""
	}
	if requestErr != nil {
		log.Printf("Error sending request to Nationalize API for name: %s, error: %v", name, requestErr)
		return ""
	}
	var countries NationalityRequest
	decodeErr := json.NewDecoder(response.Body).Decode(&countries)
	if decodeErr != nil {
		log.Printf("Error decoding JSON response from Nationalize API for name: %s, error: %v", name, decodeErr)
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
		log.Println("Country not found for name:", name)
		return ""
	}
	request := fmt.Sprintf("https://restcountries.com/v3.1/alpha/%s", code)
	response, requestErr := http.Get(request)
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()
	if requestErr != nil {
		log.Printf("Error sending request to RestCountries API for country code: %s, error: %v", code, requestErr)
		return ""
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("Error: Received non-OK status code from RestCountries API for country code: %s, status code: %d", code, response.StatusCode)
		return ""
	}

	var countriesInfo []CountryInfo
	decodeErr := json.NewDecoder(response.Body).Decode(&countriesInfo)
	if decodeErr != nil {
		log.Printf("Error decoding JSON response from RestCountries API for country code: %s, error: %v", code, decodeErr)
		return ""
	}

	return countriesInfo[0].Name.Common
}
