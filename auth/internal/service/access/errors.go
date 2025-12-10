package access

import "errors"

var (
	TokenInvalid   = errors.New("access token is invalid")
	RolesGetFailed = errors.New("failed to get accessible roles")
	Denied         = errors.New("access denied")
)
