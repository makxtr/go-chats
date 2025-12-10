package utils

import "errors"

var ErrNoMetadata = errors.New("metadata not found in context")
var ErrNoHeader = errors.New("authorization header missing")
var ErrInvalidFormat = errors.New("invalid authorization header format")
