package models

// ConvertWeight converts a weight value from one unit to another.
// Supported units: "lb" and "kg". Returns the input unchanged if units are equal or unrecognized.
func ConvertWeight(weight float64, fromUnit, toUnit string) float64 {
	if fromUnit == toUnit {
		return weight
	}
	if fromUnit == "kg" && toUnit == "lb" {
		return weight * 2.20462
	}
	if fromUnit == "lb" && toUnit == "kg" {
		return weight * 0.453592
	}
	return weight
}
