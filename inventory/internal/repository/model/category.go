package model

type PartCategory int32

// Valid PartCategory values
const (
	PartCategoryUnknown  PartCategory = 0
	PartCategoryEngine   PartCategory = 1
	PartCategoryFuel     PartCategory = 2
	PartCategoryPorthole PartCategory = 3
	PartCategoryWing     PartCategory = 4
)

// IsValid checks if the OrderStatus has a valid value
func (os PartCategory) IsValid() bool {
	switch os {
	case PartCategoryUnknown, PartCategoryEngine, PartCategoryFuel, PartCategoryPorthole, PartCategoryWing:
		return true
	default:
		return false
	}
}
