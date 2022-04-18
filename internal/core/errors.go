package core

import "errors"

var (
	ErrUserNotFound      = errors.New("user doesn't exists")
	ErrUserAlreadyExists = errors.New("user with such username already exists")

	ErrRoomNotFound = errors.New("room doesn't exists")
)
