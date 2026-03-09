package rules

import (
	"strings"
	"unicode"
)

type Violation struct {
	Message string
}

const (
	startLowerMessage     = "log message must start with a lowercase letter"
	englishOnlyMessage    = "log message must contain only English letters"
	specialSymbolsMessage = "log message must not contain special symbols or emoji"
	sensitiveDataMessage  = "log message must not contain potentially sensitive data"
)

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
		out = append(out, Violation{Message: startLowerMessage})
	}

	if !englishOnly(msg) {
		out = append(out, Violation{Message: englishOnlyMessage})
	}

	if hasSpecialSymbolsOrEmoji(msg) {
		out = append(out, Violation{Message: specialSymbolsMessage})
	}

	if containsSensitiveData(msg) {
		out = append(out, Violation{Message: sensitiveDataMessage})
	}

	return out
}

func ValidatePartialMessage(msg string) []Violation {
	if containsSensitiveData(msg) {
		return []Violation{{Message: sensitiveDataMessage}}
	}

	return nil
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