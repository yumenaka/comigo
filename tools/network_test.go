package tools

import (
	"strings"
	"testing"
)

func TestIsValidDomain(t *testing.T) {
	validLabel63 := strings.Repeat("a", 63)
	invalidLabel64 := strings.Repeat("b", 64)

	tests := []struct {
		name     string
		host     string
		expected bool
	}{
		{
			name:     "regular domain",
			host:     "example.com",
			expected: true,
		},
		{
			name:     "single label domain",
			host:     "localhost",
			expected: true,
		},
		{
			name:     "uppercase letters allowed",
			host:     "API.SERVICE",
			expected: true,
		},
		{
			name:     "hyphen inside label",
			host:     "my-service.example",
			expected: true,
		},
		{
			name:     "unicode punycode prefix",
			host:     "xn--bcher-kva.de",
			expected: true,
		},
		{
			name:     "maximum label length",
			host:     validLabel63 + ".com",
			expected: true,
		},
		{
			name:     "empty string",
			host:     "",
			expected: false,
		},
		{
			name:     "too long overall",
			host:     strings.Repeat("a", 254),
			expected: false,
		},
		{
			name:     "label too long",
			host:     invalidLabel64 + ".com",
			expected: false,
		},
		{
			name:     "leading hyphen",
			host:     "-example.com",
			expected: false,
		},
		{
			name:     "trailing hyphen",
			host:     "example-.com",
			expected: false,
		},
		{
			name:     "invalid underscore",
			host:     "exa_mple.com",
			expected: false,
		},
		{
			name:     "trailing dot",
			host:     "example.com.",
			expected: false,
		},
		{
			name:     "contains spaces",
			host:     "example .com",
			expected: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := IsValidDomain(tc.host); got != tc.expected {
				t.Errorf("IsValidDomain(%q) = %v, want %v", tc.host, got, tc.expected)
			}
		})
	}
}
