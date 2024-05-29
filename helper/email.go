package helper

import "strings"

func SanitizeEmail(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	return s
}
