package rules

import "testing"

func hasViolation(violations []Violation, msg string) bool {
	for _, v := range violations {
		if v.Message == msg {
			return true
		}
	}

	return false
}

func TestStartsWithLower(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "lowercase", input: "starting server", want: true},
		{name: "uppercase", input: "Starting server", want: false},
		{name: "leading spaces", input: "   starting server", want: true},
		{name: "empty", input: "", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := startsWithLower(tt.input); got != tt.want {
				t.Fatalf("startsWithLower(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestEnglishOnly(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "english", input: "starting server", want: true},
		{name: "english with digits", input: "server 42", want: true},
		{name: "cyrillic", input: "запуск сервера", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := englishOnly(tt.input); got != tt.want {
				t.Fatalf("englishOnly(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestHasSpecialSymbolsOrEmoji(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "plain text", input: "starting server", want: false},
		{name: "dash allowed", input: "starting-server", want: false},
		{name: "exclamation", input: "server started!", want: true},
		{name: "emoji", input: "server started 🚀", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasSpecialSymbolsOrEmoji(tt.input); got != tt.want {
				t.Fatalf("hasSpecialSymbolsOrEmoji(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestContainsSensitiveData(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "regular text", input: "starting server", want: false},
		{name: "password", input: "password leaked", want: true},
		{name: "token", input: "token expired", want: true},
		{name: "mixed case", input: "Private_Key rotated", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsSensitiveData(tt.input); got != tt.want {
				t.Fatalf("containsSensitiveData(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateMessage(t *testing.T) {
	violations := ValidateMessage("Starting password!")

	if !hasViolation(violations, startLowerMessage) {
		t.Fatalf("expected violation %q", startLowerMessage)
	}

	if !hasViolation(violations, sensitiveDataMessage) {
		t.Fatalf("expected violation %q", sensitiveDataMessage)
	}

	if !hasViolation(violations, specialSymbolsMessage) {
		t.Fatalf("expected violation %q", specialSymbolsMessage)
	}
}

func TestValidatePartialMessage(t *testing.T) {
	violations := ValidatePartialMessage("user password: ")

	if !hasViolation(violations, sensitiveDataMessage) {
		t.Fatalf("expected violation %q", sensitiveDataMessage)
	}

	violations = ValidatePartialMessage("starting server")
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}