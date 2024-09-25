package main

import "errors"

var (
	ErrUserAlreadyExist = errors.New("this email already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
)

func formatErrorToJson(err error) string {
	return `{"error": "` + err.Error() + `"}`
}
