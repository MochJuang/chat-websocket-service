package exception

import "encoding/json"

type ErrValidation Err

func Validation(err error) ErrValidation {
	return ErrValidation{
		ErrorType: TypeErrorValidation,
		ErrorCode: 400,
		Message:   err.Error(),
	}
}

func (e ErrValidation) Error() string {
	var msg string
	if IsHttpError {
		payload, _ := json.Marshal(e)
		msg = string(payload)
	}

	return msg
}
