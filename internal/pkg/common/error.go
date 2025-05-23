package pkg

import "errors"

type HttpError struct {
	Error string `json:"error"`
}

var (
	UserNotExistErr  = errors.New("user does not exist")
	SlotNotExistErr  = errors.New("appointment slot does not exist")
	SlotIsExpiredErr = errors.New("slot is expired")
)
