package utils

import (
	"encoding/json"
	"fmt"

	e "websocket-service/internal/exception"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string
	Message string
}

func Validate(s interface{}) error {
	var validationErrors []ValidationError

	validate := validator.New()

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Message: fmt.Sprintf("Validation failed on field '%s' for condition '%s'", err.Field(), err.Tag()),
			})
		}
	}

	if len(validationErrors) > 0 {
		jsonByte, err := json.Marshal(validationErrors)
		if err != nil {
			return e.Internal(err)
		}

		err = fmt.Errorf(string(jsonByte))
		return e.Validation(err)
	}

	return nil
}
