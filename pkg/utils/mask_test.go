package utils

import "testing"

func TestMaskPhone(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"+8613812345678", "+86138****5678"},
		{"13812345678", "138****5678"},
		{"+12125551234", "+12125****1234"},
		{"123", "***"},
		{"", ""},
	}

	for _, tt := range tests {
		result := MaskPhone(tt.input)
		if result != tt.expected {
			t.Errorf("MaskPhone(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestMaskName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"张三", "*三"},
		{"李四五", "*四五"},
		{"欧阳修竹", "*阳修竹"},
		{"A", "*"},
		{"", ""},
	}

	for _, tt := range tests {
		result := MaskName(tt.input)
		if result != tt.expected {
			t.Errorf("MaskName(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestMaskIDCard(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"110101199001011234", "1101**********1234"},
		{"440305200012015678", "4403**********5678"},
		{"12345678", "1234****5678"},
		{"", ""},
	}

	for _, tt := range tests {
		result := MaskIDCard(tt.input)
		if result != tt.expected {
			t.Errorf("MaskIDCard(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}
