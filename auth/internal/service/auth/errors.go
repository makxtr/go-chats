package auth

import "errors"

var (
	TokenGenFailed       = errors.New("token generation failed")
	LoginFailed          = errors.New("login failed")
	RefreshTokenInvalid  = errors.New("refresh token is invalid")
)
