package core

import "errors"

var ErrResourceExhausted = errors.New("request limit exhausted")
var ErrInvalidType = errors.New("limit must be a number")
var ErrFailedToProcessLimit = errors.New("failed to process limit")
