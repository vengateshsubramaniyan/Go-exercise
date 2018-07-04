package main

import (
	"testing"
)

type testCase struct {
	input string
	want  string
}

func TestNormalize(t *testing.T) {
	tests := []testCase{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"(123)456-7892", "1234567892"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			actual := normalize(tc.input)
			if tc.want != actual {
				t.Errorf("got %s, but want %s\n", actual, tc.want)
			}
		})
	}
}
