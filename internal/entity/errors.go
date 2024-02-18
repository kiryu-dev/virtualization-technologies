package entity

import (
	"errors"
)

var (
	ErrUserNotFound        = errors.New("specified user is not found")
	ErrEmailIsAlreadyTaken = errors.New("specified email is already taken")
)
