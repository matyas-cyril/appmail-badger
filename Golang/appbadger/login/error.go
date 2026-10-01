package login

import (
	"errors"
	"fmt"
)

var (
	ErrLoginInvalidLength = fmt.Errorf("login must be between %d and %d characters", LOGIN_MIN_SIZE, LOGIN_MAX_SIZE)
	ErrLoginInvalidFormat = errors.New("login contains characters that are not allowed")

	ErrPasswordInvalidLength = fmt.Errorf("password must be between %d and %d characters", PASSWD_MIN_SIZE, PASSWD_MAX_SIZE)
	ErrPasswordInvalidFormat = errors.New("password contains characters that are not allowed")
)
