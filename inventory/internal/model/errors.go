package model

import (
	"errors"

	sharedErrors "github.com/andredubov/rocket-factory/shared/pkg/errors"
)

var (
	ErrPartNotFound      = sharedErrors.NewNotFoundError(errors.New("part not found"))
	ErrInvalidUUID       = sharedErrors.NewInvalidArgumentError(errors.New("invalid uuid"))
	ErrPartAlreadyExists = sharedErrors.NewAlreadyExistsError(errors.New("part already exists"))
)

// import "errors"

// Error definitions for the inventory repository.
// These are base errors that can be wrapped with additional context.
// var (
// 	ErrPartAlreadyExists = errors.New("part already exists")
// 	ErrPartNotFound      = errors.New("part not found")
// )
