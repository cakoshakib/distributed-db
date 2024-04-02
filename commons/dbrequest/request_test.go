package network

import (
	"assert"
	"testing"
)

func TestRequest_StringToOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       string
		expected operation
	}{
		{
			name:     "happy case",
			in:       "get",
			expected: Get,
		},
		{
			name:     "unspecified case",
			in:       "notAnOp",
			expected: Unspecified,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := StringToOperation(tt.in)
			assert.Equal(out, tt.expected)
		})
	}
}
