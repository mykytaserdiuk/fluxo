package fluxo

import "errors"

var (
	ErrNotAFunc     = errors.New("fn params is not a function")
	ErrFuncIsNil    = errors.New("fn params is nil")
	ErrFnIsNotValid = errors.New("fn params is not a valid")
	ErrNoHandlers   = errors.New("id hasn`t handlers")
)
