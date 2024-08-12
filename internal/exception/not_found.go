package exception

import "encoding/json"

type ErrNotFound Err

func NotFound(message string) ErrNotFound {
	return ErrNotFound{
		ErrorType: TypeErrorNotFound,
		ErrorCode: 404,
		Message:   message,
	}
}

func (e ErrNotFound) Error() string {
	var msg string
	if IsHttpError {
		payload, _ := json.Marshal(e)
		msg = string(payload)
	}

	return msg
}
