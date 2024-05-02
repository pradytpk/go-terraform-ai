package utils

import "strings"

// ToPtr converts type T to a *T as a convenience.
//
//	@param i
//	@return *T
func ToPtr[T any](i T) *T {
	return &i
}

// RemoveBlankLinesFromString remove the blaonk lines from the code
//
//	@param input
//	@return string
func RemoveBlankLinesFromString(input string) string {
	return strings.TrimLeft(input, "\n\r \t")
}
