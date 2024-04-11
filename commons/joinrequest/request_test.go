package joinrequest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequest_NewRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       string
		expected JoinRequest
	}{
		{
			name: "happy case",
			in:   "join node-12 127.0.0.1;",
			expected: JoinRequest{
				NodeID:  "node-12",
				Address: "127.0.0.1",
			},
		},
		{
			name:     "bad case",
			in:       "join ;",
			expected: JoinRequest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := NewRequest(tt.in)
			assert.Equal(t, out, tt.expected)
		})
	}
}

func TestRequest_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       JoinRequest
		expected bool
	}{
		{
			name: "happy case",
			in: JoinRequest{
				NodeID:  "node-12",
				Address: "127.0.0.1",
			},
			expected: true,
		},
		{
			name: "no node id",
			in: JoinRequest{
				NodeID:  "",
				Address: "127.0.0.1",
			},
			expected: false,
		},
		{
			name: "no address",
			in: JoinRequest{
				NodeID:  "node-12",
				Address: "",
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.in.Validate()
			assert.Equal(t, out, tt.expected)
		})
	}
}
