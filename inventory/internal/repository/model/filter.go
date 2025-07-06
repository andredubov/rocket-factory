package model

// PartFilter contains criteria for searching inventory parts
// All fields are optional - empty slices mean no filtering on that field
// Multiple values within a field are treated as OR conditions
// Different fields are combined with AND logic
type PartFilter struct {
	UUIDs                 []string   // Filter by part UUIDs
	Names                 []string   // Filter by part names
	Categories            []Category // Filter by categories
	ManufacturerCountries []string   // Filter by manufacturer countries
	Tags                  []string   // Filter by part tags
}
