package model

import (
	"errors"
	"fmt"

	sharedErrors "github.com/andredubov/rocket-factory/shared/pkg/errors"
)

var (
	ErrPartNotFound      = sharedErrors.NewNotFoundError(errors.New("part not found"))
	ErrInvalidUUID       = sharedErrors.NewInvalidArgumentError(errors.New("invalid uuid"))
	ErrPartAlreadyExists = sharedErrors.NewAlreadyExistsError(errors.New("part already exists"))
)

// ErrPartWithUUIDNotFound constructs a new error indicating a part with the specified UUID was not found.
func ErrPartWithUUIDNotFound(uuid string) error {
	return fmt.Errorf("part with UUID %s not found: %w", uuid, ErrPartNotFound)
}

// ErrPartWithUUIDExists constructs a new error indicating a part with the specified UUID already exists.
func ErrPartWithUUIDExists(uuid string) error {
	return fmt.Errorf("part with UUID %s already exists: %w", uuid, ErrPartAlreadyExists)
}

// import "errors"

// Error definitions for the inventory repository.
// These are base errors that can be wrapped with additional context.
// var (
// 	ErrPartAlreadyExists = errors.New("part already exists")
// 	ErrPartNotFound      = errors.New("part not found")
// )
