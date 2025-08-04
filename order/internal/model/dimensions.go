package model

// Dimensions contains physical measurement data for inventory parts
// All values are stored in metric units (centimeters and kilograms)
type Dimensions struct {
	Length float64 // Length in centimeters
	Width  float64 // Width in centimeters
	Height float64 // Height in centimeters
	Weight float64 // Weight in kilograms
}
