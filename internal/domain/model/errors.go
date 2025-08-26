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
	ErrWeakPassword         = errors.New("weak Password")
	ErrBadFormatEmail       = errors.New("email does not match the required format")
	ErrUserNotFound         = errors.New("user not found")

	ErrInvalidChannelName      = errors.New("invalid Channel Name")
	ErrBadFormatChannelName    = errors.New("Channel Name does not match the required format")
	ErrAlreadyExistChannelName = errors.New("Channel Name already exists")
	ErrInvalidDisplayName      = errors.New("invalid Channel Display Name")
	ErrChannelNotFound         = errors.New("channel not found")

	ErrInvalidMessageContent = errors.New("invalid Message Content")
	ErrMessageNotFound       = errors.New("message not found")
	ErrInvalidRequestLimit   = errors.New("invalid request limit")
	ErrInvalidTimeRange      = errors.New("invalid time range")
	ErrMessageAlreadyPinned  = errors.New("message already pinned")
	ErrMessageNotPinned      = errors.New("message not pinned")
)
