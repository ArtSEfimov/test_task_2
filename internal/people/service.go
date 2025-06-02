package people

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

func enrichPerson(person *Person) {
	person.Age = getAge(person.Name)
	person.Gender = getGender(person.Name)
}
