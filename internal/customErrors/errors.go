package customErrors

import "errors"

var ErrInvalidAmount = errors.New("сумма должна быть положительной")
var ErrNoMoney = errors.New("недостаточно средств")
var ErrUserNotFound = errors.New("пользователь не найден")
