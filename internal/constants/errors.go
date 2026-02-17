package constants

import "errors"

var (
	ErrUsernameOrEmailAlreadyExists = errors.New("Username or Email Already Exists!")
	ErrDataNotFound                 = errors.New("Data Not Found")
	ErrInvalidToken                 = errors.New("Invalid Token")
	ErrInvalidPassword              = errors.New("Invalid Password")
)
