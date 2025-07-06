package model

// Category represents a classification group for inventory parts
// Used to organize and filter parts in the inventory system
type Category struct {
	ID   string // Unique identifier for the category
	Name string // Human-readable name of the category
}
