package converter

import "testing"

func TestParseQuality(t *testing.T) {
	tests := []struct {
		input    string
		expected Quality
	}{
		{"screen", Screen},
		{"printer", Printer},
		{"default", Default},
		{"invalid", Default},
		{"", Default},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseQuality(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestQualityString(t *testing.T) {
	tests := []struct {
		quality  Quality
		expected string
	}{
		{Screen, "screen"},
		{Printer, "printer"},
		{Default, "default"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.quality.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
