package services

import "errors"

var (
	ErrUnknown          = errors.New("unknown")
	ErrNotFound         = errors.New("not found")
	ErrUpdateFailed     = errors.New("update failed")
	ErrConflict         = errors.New("conflict")
	ErrSameCourseCommon = errors.New("same course common")
	ErrWrongPassword    = errors.New("wrong password")
	ErrInvalidParams    = errors.New("invalid params")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrQuotaExceeded    = errors.New("quota exceeded")
	ErrCannotOperate    = errors.New("cannot operate")
)
