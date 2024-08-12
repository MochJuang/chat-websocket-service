package exception

import "encoding/json"

type ErrUnauthorized Err

func Unauthorized(err error) ErrUnauthorized {
	return ErrUnauthorized{
		ErrorType: TypeErrorUnauthorized,
		ErrorCode: 401,
		Message:   err.Error(),
	}
}

func (e ErrUnauthorized) Error() string {
	var msg string
	if IsHttpError {
		payload, _ := json.Marshal(e)
		msg = string(payload)
	}

	return msg
}
