package config

import "fmt"

type Config struct {
	CheckLowercase      bool
	CheckEnglishOnly    bool
	CheckSpecialSymbols bool
	CheckSensitiveData  bool
	AllowDash           bool
	SensitiveKeywords   []string
}

func Default() Config {
	return Config{
		CheckLowercase:      true,
		CheckEnglishOnly:    true,
		CheckSpecialSymbols: true,
		CheckSensitiveData:  true,
		AllowDash:           true,
		SensitiveKeywords: []string{
			"password",
			"passwd",
			"token",
			"api_key",
			"apikey",
			"secret",
			"access_key",
			"private_key",
		},
	}
}

func FromAny(settings any) (Config, error) {
	cfg := Default()
	if settings == nil {
		return cfg, nil
	}

	m, ok := settings.(map[string]any)
	if !ok {
		return cfg, fmt.Errorf("unexpected settings type %T", settings)
	}

	if v, ok := boolValue(m, "check-lowercase"); ok {
		cfg.CheckLowercase = v
	}
	if v, ok := boolValue(m, "check-english-only"); ok {
		cfg.CheckEnglishOnly = v
	}
	if v, ok := boolValue(m, "check-special-symbols"); ok {
		cfg.CheckSpecialSymbols = v
	}
	if v, ok := boolValue(m, "check-sensitive-data"); ok {
		cfg.CheckSensitiveData = v
	}
	if v, ok := boolValue(m, "allow-dash"); ok {
		cfg.AllowDash = v
	}
	if v, ok := stringSliceValue(m, "sensitive-keywords"); ok && len(v) > 0 {
		cfg.SensitiveKeywords = v
	}

	return cfg, nil
}

func boolValue(m map[string]any, key string) (bool, bool) {
	v, ok := m[key]
	if !ok {
		return false, false
	}

	b, ok := v.(bool)
	return b, ok
}

func stringSliceValue(m map[string]any, key string) ([]string, bool) {
	v, ok := m[key]
	if !ok {
		return nil, false
	}

	switch raw := v.(type) {
	case []string:
		return raw, true
	case []any:
		out := make([]string, 0, len(raw))
		for _, item := range raw {
			s, ok := item.(string)
			if !ok {
				return nil, false
			}
			out = append(out, s)
		}
		return out, true
	default:
		return nil, false
	}
}