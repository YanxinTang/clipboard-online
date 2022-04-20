package utils

import (
	"path/filepath"
	"testing"
)

func TestAppendOrderToFilename(t *testing.T) {
	tcs := []struct {
		input string
		want  string
	}{
		{"bar.png", "bar(1).png"},
		{"/foo/bar.png", "/foo/bar(1).png"},
		{"bar(1).png", "bar(2).png"},
		{"/foo/bar(1).png", "/foo/bar(2).png"},
		{"(1).png", "(2).png"},
		{"/foo/(1).png", "/foo/(2).png"},
		{"(1)", "(2)"},
		{"/foo/(1)", "/foo/(2)"},
	}

	for _, tc := range tcs {
		got := AppendOrderToFilename(tc.input)
		if filepath.Clean(tc.want) != got {
			t.Errorf("AppendOrderToFilename(%s) = %s, want %s", tc.input, got, tc.want)
		}
	}
}
