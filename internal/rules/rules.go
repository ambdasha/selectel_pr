package rules

import (
	"strings"
	"unicode"
)

type Violation struct {
	Message string
}

var sensitiveKeywords = []string{
	"password",
	"passwd",
	"token",
	"api_key",
	"apikey",
	"secret",
	"access_key",
	"private_key",
}

func ValidateMessage(msg string) []Violation {
	var out []Violation

	if !startsWithLower(msg) {
		out = append(out, Violation{
			Message: "log message must start with a lowercase letter",
		})
	}

	if !englishOnly(msg) {
		out = append(out, Violation{
			Message: "log message must contain only English letters",
		})
	}

	if hasSpecialSymbolsOrEmoji(msg) {
		out = append(out, Violation{
			Message: "log message must not contain special symbols or emoji",
		})
	}

	if containsSensitiveData(msg) {
		out = append(out, Violation{
			Message: "log message must not contain potentially sensitive data",
		})
	}

	return out
}

func startsWithLower(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}

	for _, r := range s {
		if unicode.IsLetter(r) {
			return unicode.IsLower(r)
		}
	}

	return true
}

func englishOnly(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func hasSpecialSymbolsOrEmoji(s string) bool {
	for _, r := range s {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), unicode.IsSpace(r):
			continue
		case r == '-':
			continue
		default:
			return true
		}
	}
	return false
}

func containsSensitiveData(s string) bool {
	low := strings.ToLower(s)
	for _, kw := range sensitiveKeywords {
		if strings.Contains(low, kw) {
			return true
		}
	}
	return false
}