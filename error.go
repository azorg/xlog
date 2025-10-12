// File: "error.go"

package xlog

import "errors"

// Ошибка: "ротация файла журнала не предусмотрена конфигурацией"
var ErrNotRotatable = errors.New("logger is not rotatable")

// Ошибка: не передан хендлер (nil)
var ErrNilHandler = errors.New("handler is nil")

// EOF: "error.go"
