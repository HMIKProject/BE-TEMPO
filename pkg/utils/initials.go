package utils

import (
	"strings"
	"unicode"
)

// GenerateInitials menghasilkan 1 atau 2 huruf inisial kapital dari nama
func GenerateInitials(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "AH"
	}

	parts := strings.Fields(name)
	if len(parts) >= 2 {
		r1 := []rune(parts[0])
		r2 := []rune(parts[1])
		if len(r1) > 0 && len(r2) > 0 {
			return string(unicode.ToUpper(r1[0])) + string(unicode.ToUpper(r2[0]))
		}
	}

	runes := []rune(name)
	if len(runes) >= 2 {
		return string(unicode.ToUpper(runes[0])) + string(unicode.ToUpper(runes[1]))
	}
	if len(runes) == 1 {
		return string(unicode.ToUpper(runes[0]))
	}

	return "AH"
}
