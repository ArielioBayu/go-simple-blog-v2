package constants

import "errors"

var (
	ErrUsernameOrEmailAlreadyExists = errors.New("Username or Email Already Exists!")
	ErrDataNotFound                 = errors.New("Data Not Found")
)
