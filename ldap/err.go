package ldap

import "errors"

var (
	ErrUserNotFound           = errors.New("Пользователь не найден")
	ErrUserFoundMoreThanOne   = errors.New("Найденно пользователей больше чем один")
	ErrUserSAMAccountRequired = errors.New("SAMAccount обезателен")
)
