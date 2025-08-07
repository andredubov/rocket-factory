package model

// Dimensions contains physical measurement data for inventory parts
// All values are stored in metric units (centimeters and kilograms)
type Dimensions struct {
	Length float64 `bson:"length"` // Length in centimeters
	Width  float64 `bson:"width"`  // Width in centimeters
	Height float64 `bson:"height"` // Height in centimeters
	Weight float64 `bson:"weight"` // Weight in kilograms
}
