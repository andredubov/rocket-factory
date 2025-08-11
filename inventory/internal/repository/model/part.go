package model

import (
	"time"
)

// Value represents a typed value for metadata fields
// Uses pointers to distinguish between zero values and unset values
type Value struct {
	StringValue *string  `bson:"stringValue,omitempty"` // String type metadata value
	Int64Value  *int64   `bson:"int64Value,omitempty"`  // Integer type metadata value
	DoubleValue *float64 `bson:"doubleValue,omitempty"` // Floating-point type metadata value
	BoolValue   *bool    `bson:"boolValue,omitempty"`   // Boolean type metadata value
}

// Part represents an inventory item with all associated data
// This is the core domain entity for the inventory system
type Part struct {
	Uuid          string           `bson:"_id"`                 // Unique identifier
	Name          string           `bson:"name"`                // Display name
	Description   string           `bson:"description"`         // Detailed description
	Price         float64          `bson:"price"`               // Unit price in base currency
	StockQuantity int64            `bson:"stockQuantity"`       // Current inventory count
	Category      PartCategory     `bson:"category"`            // Classification group
	Dimensions    Dimensions       `bson:"dimensions"`          // Physical measurements
	Manufacturer  Manufacturer     `bson:"manufacturer"`        // Producer information
	Tags          []string         `bson:"tags,omitempty"`      // Searchable keywords
	Metadata      map[string]Value `bson:"metadata,omitempty"`  // Flexible key-value storage
	CreatedAt     time.Time        `bson:"createdAt"`           // When part was added to system
	UpdatedAt     time.Time        `bson:"updatedAt,omitempty"` // When part was last modified
}
