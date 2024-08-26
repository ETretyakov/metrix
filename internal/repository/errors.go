package repository

import "errors"

// ErrNotFound - the variable that stores specific NotFound error type.
var ErrNotFound = errors.New("entity not found")
