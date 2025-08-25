package model

import "errors"

var (
	ErrInvalidUserName      = errors.New("invalid User Name")
	ErrBadFormatUserName    = errors.New("User Name does not match the required format")
	ErrAlreadyExistUserName = errors.New("User Name already exists")
	ErrInvalidNickname      = errors.New("invalid Nickname")
)
