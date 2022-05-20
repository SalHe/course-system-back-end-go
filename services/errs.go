package services

import "errors"

var (
	ErrUnknown       = errors.New("unknown")
	ErrNotFound      = errors.New("not found")
	ErrUpdateFailed  = errors.New("update failed")
	ErrConflict      = errors.New("conflict")
	ErrWrongPassword = errors.New("wrong password")
	ErrInvalidParams = errors.New("invalid params")
)
