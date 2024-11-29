package db2go

import "strings"

// Camelize converts a string with underscores into camelCase.
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
