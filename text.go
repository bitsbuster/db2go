package db2go

import "strings"

// Camelize converts a snake_case string into a camelCase or PascalCase string.
//
// This function transforms an input string from snake_case to camelCase or PascalCase,
// depending on the `capitalised` flag. Each underscore-separated word in the input string
// is capitalized, and underscores are removed.
//
// Parameters:
//   - input: string - The snake_case string to convert.
//   - capitalised: bool - A flag indicating whether the first letter of the resulting string
//     should be capitalized (PascalCase) or lowercase (camelCase).
//
// Returns:
//   - string: The camelCase or PascalCase representation of the input string.
//
// Example Usage:
//   - Camelize("example_input", false) -> "exampleInput"
//   - Camelize("example_input", true)  -> "ExampleInput"
//
// Notes:
//   - If the input string is empty or contains no underscores, it is returned unchanged.
//   - The function assumes the input string is in valid snake_case format.
func Camelize(input string, capitalised bool) string {
	words := strings.Split(input, "_")
	for i := range words {

		if i > 0 && len(words[i]) > 0 {
			words[i] = strings.ToUpper(string(words[i][0])) + words[i][1:]
		} else if capitalised && i == 0 && len(words[i]) > 0 {
			words[i] = strings.ToUpper(string(words[i][0])) + words[i][1:]
		}
	}
	return strings.Join(words, "")
}
