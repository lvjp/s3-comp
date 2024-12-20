package client

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJoinURIPath(t *testing.T) {
	testCases := []struct{ a, b, expected string }{
		{"", "", "/"},
		{"", "/", "/"},
		{"", "abc", "/abc"},
		{"", "/abc", "/abc"},

		{"/", "", "/"},
		{"/", "/", "/"},
		{"/", "abc", "/abc"},
		{"/", "/abc", "/abc"},

		{"123", "", "/123"},
		{"123", "/", "/123"},
		{"123", "abc", "/123/abc"},
		{"123", "/abc", "/123/abc"},

		{"/123", "", "/123"},
		{"/123", "/", "/123"},
		{"/123", "abc", "/123/abc"},
		{"/123", "/abc", "/123/abc"},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := joinURIPath(tc.a, tc.b)
			require.Equal(t, tc.expected, actual)
		})
	}
}
