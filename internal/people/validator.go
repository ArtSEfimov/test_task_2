package people

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func IsValid(payload *Request) (errors []string) {
	peopleValidator := validator.New()
	validateErr := peopleValidator.Struct(&payload)

	if validateErr != nil {
		for _, e := range validateErr.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field: %s, Tag: %s", e.Field(), e.Tag()))
		}

		return
	}
	return
}
