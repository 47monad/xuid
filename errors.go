package xuid

import "errors"

var (
	ErrInvalidUUIDString = errors.New("UUID string is invalid")
	ErrParse             = errors.New("XUID string cannot be parsed")
)
