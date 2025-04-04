package core

import "errors"

var ErrResourceExhausted = errors.New("request limit exhausted")
var ErrAlreadyExists = errors.New("resource or task already exists")
