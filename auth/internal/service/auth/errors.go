package auth

import "errors"

var (
	TokenGenFailed = errors.New("token generation failed")
	LoginFailed    = errors.New("login failed")
)
