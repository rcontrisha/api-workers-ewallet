package msgerror

import "errors"

var (
	ErrNoRows       = errors.New("no row affected")
	ErrUserNotFound = errors.New("user not found")
	ErrFavNotFound  = errors.New("user not found")
)
