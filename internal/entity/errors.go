package entity

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidUserName   = errors.New("invalid user name")
	ErrInvalidUserEmail  = errors.New("invalid user email")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserID     = errors.New("invalid user ID")
)
