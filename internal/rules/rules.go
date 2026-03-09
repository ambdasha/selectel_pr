package rules

import (
	"strings"
	"unicode"

	"github.com/ambdasha/logmsglint/internal/config"
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

func ValidateMessage(msg string, cfg config.Config) []Violation {
	var out []Violation

	if cfg.CheckLowercase && !startsWithLower(msg) {
		out = append(out, Violation{Message: startLowerMessage})
	}

	if cfg.CheckEnglishOnly && !englishOnly(msg) {
		out = append(out, Violation{Message: englishOnlyMessage})
	}

	if cfg.CheckSpecialSymbols && hasSpecialSymbolsOrEmoji(msg, cfg.AllowDash) {
		out = append(out, Violation{Message: specialSymbolsMessage})
	}

	if cfg.CheckSensitiveData && containsSensitiveData(msg, cfg.SensitiveKeywords) {
		out = append(out, Violation{Message: sensitiveDataMessage})
	}

	return out
}

func ValidatePartialMessage(msg string, cfg config.Config) []Violation {
	if cfg.CheckSensitiveData && containsSensitiveData(msg, cfg.SensitiveKeywords) {
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

func hasSpecialSymbolsOrEmoji(s string, allowDash bool) bool {
	for _, r := range s {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), unicode.IsSpace(r):
			continue
		case allowDash && r == '-':
			continue
		default:
			return true
		}
	}

	return false
}

func containsSensitiveData(s string, sensitiveKeywords []string) bool {
	low := strings.ToLower(s)
	for _, kw := range sensitiveKeywords {
		if strings.Contains(low, strings.ToLower(kw)) {
			return true
		}
	}

	return false
}