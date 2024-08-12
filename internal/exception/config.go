package exception

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
)

const IsHttpError = true
const (
	TypeErrorNotFound     = "NotFound"
	TypeErrorValidation   = "Validation"
	TypeErrorUnauthorized = "Unauthorized"
	TypeErrorInternal     = "Internal"
)

type Err struct {
	ErrorType string
	ErrorCode int
	Message   string
}

func (e Err) Error() string {
	var msg string
	if IsHttpError {
		payload, _ := json.Marshal(e)
		msg = string(payload)
	}

	return msg
}

func Convert(err error) (Err, error) {
	var e Err
	newErr := json.Unmarshal([]byte(err.Error()), &e)
	if newErr != nil {
		return e, newErr
	}

	return e, nil
}

func HandleHttpErrorFiber(c *fiber.Ctx, err error) error {
	var msg, _ = Convert(err)
	log.Println(msg)
	return c.Status(msg.ErrorCode).JSON(fiber.Map{"errors": msg.Message})
}
