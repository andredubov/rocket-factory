package model

import "errors"

// Error definitions for the inventory repository.
// These are base errors that can be wrapped with additional context.
var (
	ErrPartAlreadyExists = errors.New("part already exists")
	ErrPartNotFound      = errors.New("part not found")
)
