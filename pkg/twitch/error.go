package twitch

import (
	"errors"
	"fmt"
)

var ErrUnauthorized = errors.New("unauthorized")

type BadRequestError struct {
	Message string
}

func (b BadRequestError) Error() string {
	return fmt.Sprintf("bad request: %q", b.Message)
}
