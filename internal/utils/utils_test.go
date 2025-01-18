package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRightPadTrim(t *testing.T) {
	inputStr := "hello"
	tests := []struct {
		input    string
		length   int
		dots     bool
		expected string
	}{
		{inputStr, 10, false, "hello     "},
		{inputStr, 10, true, "hello     "},
		{inputStr, 5, false, "hello"},
		{inputStr, 5, true, "hello"},
		{inputStr, 4, false, "hell"},
		{inputStr, 4, true, "h..."},
		{inputStr, 3, true, "hel"},
		{inputStr, 3, false, "hel"},
		{inputStr, 2, true, "he"},
		{inputStr, 2, false, "he"},
		{inputStr, 1, false, "h"},
		{inputStr, 0, false, ""},
	}

	for _, tc := range tests {
		got := RightPadTrim(tc.input, tc.length, tc.dots)
		assert.Equal(t, tc.expected, got, "length: %d, dots: %v", tc.length, tc.dots)
	}
}

func TestTrim(t *testing.T) {
	inputStr := "hello"
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{inputStr, 5, "hello"},
		{inputStr, 6, "hello"},
		{inputStr, 4, "h..."},
		{inputStr, 3, "hel"},
		{inputStr, 2, "he"},
		{inputStr, 1, "h"},
		{inputStr, 0, ""},
	}

	for _, tc := range tests {
		got := Trim(tc.input, tc.length)
		assert.Equal(t, tc.expected, got, "input: %s, length: %d", tc.input, tc.length)
	}
}

func TestTrimWithMoreLinesIndicator(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello", 4, "h..."},
		{"hello", 3, "hel"},
		{"hello\nworld", 10, "hello ~"},
		{"hello\nworld", 7, "hello ~"},
		{"hello\nworld", 5, "hello"},
		{"hello\nworld", 4, "h..."},
		{"hello\nworld", 3, "hel"},
	}

	for _, tc := range tests {
		got := TrimWithMoreLinesIndicator(tc.input, tc.length)
		assert.Equal(t, tc.expected, got, "input: %s, length: %d", tc.input, tc.length)
	}
}

func TestRightPadTrimWithMoreLinesIndicator(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{"hello", 10, "hello     "},
		{"hello", 5, "hello"},
		{"hello", 4, "h..."},
		{"hello", 3, "hel"},
		{"hello\nworld", 10, "hello ~   "},
		{"hello\nworld", 7, "hello ~"},
		{"hello\nworld", 5, "hello"},
		{"hello\nworld", 4, "h..."},
		{"hello\nworld", 3, "hel"},
	}

	for _, tc := range tests {
		got := RightPadTrimWithMoreLinesIndicator(tc.input, tc.length)
		assert.Equal(t, tc.expected, got, "input: %s, length: %d", tc.input, tc.length)
	}
}
