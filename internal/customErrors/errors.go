package customErrors

import "errors"

var (
	ErrInvalidAmount = errors.New("сумма должна быть положительной")
	ErrNoMoney       = errors.New("недостаточно средств")
	ErrUserNotFound  = errors.New("пользователь не найден")
)
