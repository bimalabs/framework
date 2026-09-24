package strcase

import "strings"

func ToDelimited(s string, delimiter byte) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	var result strings.Builder
	result.Grow(len(s))
	for i := range len(s) {
		current := s[i]
		upper := current >= 'A' && current <= 'Z'
		lower := current >= 'a' && current <= 'z'
		digit := current >= '0' && current <= '9'
		if upper {
			current += 'a' - 'A'
		}

		if i+1 < len(s) {
			next := s[i+1]
			nextUpper := next >= 'A' && next <= 'Z'
			nextLower := next >= 'a' && next <= 'z'
			nextDigit := next >= '0' && next <= '9'
			if upper && (nextLower || nextDigit) || lower && (nextUpper || nextDigit) || digit && (nextUpper || nextLower) {
				if upper && nextLower && i > 0 && s[i-1] >= 'A' && s[i-1] <= 'Z' {
					result.WriteByte(delimiter)
				}

				result.WriteByte(current)
				if lower || digit || nextDigit {
					result.WriteByte(delimiter)
				}

				continue
			}
		}

		switch current {
		case ' ', '_', '-', '.':
			result.WriteByte(delimiter)
		default:
			result.WriteByte(current)
		}
	}

	return result.String()
}
