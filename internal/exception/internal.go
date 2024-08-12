package exception

import "encoding/json"

type ErrInternal Err

func Internal(err error) ErrInternal {
	return ErrInternal{
		ErrorCode: 500,
		ErrorType: TypeErrorInternal,
		Message:   err.Error(),
	}
}

func (e ErrInternal) Error() string {
	var msg string
	payload, _ := json.Marshal(e)
	msg = string(payload)

	return msg
}
