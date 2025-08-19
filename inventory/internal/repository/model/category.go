package model

import "errors"

type PartCategory int32

const (
	PartCategoryUnknown  PartCategory = 0
	PartCategoryEngine   PartCategory = 1
	PartCategoryFuel     PartCategory = 2
	PartCategoryPorthole PartCategory = 3
	PartCategoryWing     PartCategory = 4
)

var ErrUnknownPartCategory = errors.New("invalid part category")

// NewPartCategory creates a new PartCategory from a numeric value.
func NewPartCategory(number int32) (PartCategory, error) {
	switch number {
	case 0:
		return PartCategoryUnknown, nil
	case 1:
		return PartCategoryEngine, nil
	case 2:
		return PartCategoryFuel, nil
	case 3:
		return PartCategoryPorthole, nil
	case 4:
		return PartCategoryWing, nil
	default:
		return PartCategoryUnknown, ErrUnknownPartCategory
	}
}
