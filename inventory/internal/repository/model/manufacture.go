package model

// Manufacturer contains information about the part producer
type Manufacturer struct {
	Name    string `bson:"name"`              // Company name
	Country string `bson:"country"`           // Country code
	Website string `bson:"website,omitempty"` // Optional field
}
