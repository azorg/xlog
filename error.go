// File: "error.go"

package xlog

import "errors"

// Ошибка: "ротация файла журнала не предусмотрена конфигурацией"
var ErrNotRotatable = errors.New("logger is not rotatable")

// EOF: "error.go"
