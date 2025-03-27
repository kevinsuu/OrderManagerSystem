package repository

import "errors"

var (
	ErrNotFound     = errors.New("record not found")
	ErrTokenExpired = errors.New("token expired")
)
