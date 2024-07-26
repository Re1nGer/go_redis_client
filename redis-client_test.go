package redisclient

import (
	"bytes"
	"testing"
)

// Helper function to compare byte slices
func bytesEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}

// Test the encodeCommand function
func TestEncodeCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []byte
	}{
		{
			name:     "Simple SET command",
			args:     []string{"SET", "key", "value"},
			expected: []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"),
		},
		{
			name:     "GET command",
			args:     []string{"GET", "key"},
			expected: []byte("*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"),
		},
		{
			name:     "Complex SADD command",
			args:     []string{"SADD", "myset", "value1", "value2", "value3"},
			expected: []byte("*5\r\n$4\r\nSADD\r\n$5\r\nmyset\r\n$6\r\nvalue1\r\n$6\r\nvalue2\r\n$6\r\nvalue3\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encodeCommand(tt.args)
			if !bytesEqual(result, tt.expected) {
				t.Errorf("encodeCommand(%v) = %v, want %v", tt.args, result, tt.expected)
			}
		})
	}
}
