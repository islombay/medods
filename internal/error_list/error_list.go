package error_list

import "errors"

var (
	NotFound     = errors.New("not found")
	Unauthorized = errors.New("unauthorized")
	Forbidden    = errors.New("forbidden")
	TokenExpired = errors.New("token expired")
)
