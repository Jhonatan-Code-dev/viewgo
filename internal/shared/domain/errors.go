// Package domain provides shared kernel domain primitives and sentinel errors.
package domain

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrInvalidInput  = errors.New("invalid input parameter")
	ErrInternalError = errors.New("internal server error")
)
