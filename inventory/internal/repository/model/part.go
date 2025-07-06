package model

import (
	"time"
)

// Value represents a typed value for metadata fields
// Uses pointers to distinguish between zero values and unset values
type Value struct {
	StringValue *string  // String type metadata value
	Int64Value  *int64   // Integer type metadata value
	DoubleValue *float64 // Floating-point type metadata value
	BoolValue   *bool    // Boolean type metadata value
}

// Part represents an inventory item with all associated data
// This is the core domain entity for the inventory system
type Part struct {
	Uuid          string           // Unique identifier
	Name          string           // Display name
	Description   string           // Detailed description
	Price         float64          // Unit price in base currency
	StockQuantity int64            // Current inventory count
	Category      Category         // Classification group
	Dimensions    Dimensions       // Physical measurements
	Manufacturer  Manufacturer     // Producer information
	Tags          []string         // Searchable keywords
	Metadata      map[string]Value // Flexible key-value storage
	CreatedAt     time.Time        // When part was added to system
	UpdatedAt     time.Time        // When part was last modified
}
