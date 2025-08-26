package model

import "errors"

var (
	ErrNilUUID        = errors.New("nil UUID")
	ErrInvalidUUID    = errors.New("invalid UUID")
	ErrNothingChanged = errors.New("nothing changed")

	ErrInvalidUserName      = errors.New("invalid User Name")
	ErrBadFormatUserName    = errors.New("User Name does not match the required format")
	ErrAlreadyExistUserName = errors.New("User Name already exists")
	ErrInvalidNickname      = errors.New("invalid Nickname")
)
