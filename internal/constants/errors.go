package constants

import "errors"

var (

	//	Error
	ErrUsernameOrEmailAlreadyExists = errors.New("Username or Email Already Exists!")
	ErrMissingToken                 = errors.New("Missing Token")
	ErrUnauthorized                 = errors.New("Unauthorized")

	// Not Found
	ErrDataNotFound = errors.New("Data Not Found")
	ErrPostNotFound = errors.New("Error Post Not Found")

	//	Invalid Token
	ErrInvalidToken    = errors.New("invalid token")
	ErrInvalidPassword = errors.New("invalid password")
	ErrTokenExpired    = errors.New("refresh token has expired")
)
