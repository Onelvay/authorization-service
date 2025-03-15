package grant

import "errors"

var (
	ErrUserExist = errors.New("Пользователь с таким логином существует")
)
