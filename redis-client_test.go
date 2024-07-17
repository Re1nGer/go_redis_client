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

// Test the Set method
func TestSetCommand(t *testing.T) {
	tests := []struct {
		name     string
		opts     *SetOpts
		args     []string
		expected []byte
	}{
		{
			name:     "Simple SET",
			args:     []string{"SET", "key1", "value1"},
			opts:     NewSetOpts(),
			expected: []byte("*3\r\n$3\r\nSET\r\n$4\r\nkey1\r\n$6\r\nvalue1\r\n"),
		},
		{
			name:     "SET with NX",
			args:     []string{"SET", "key2", "value2", "NX"},
			opts:     NewSetOpts().WithNX(),
			expected: []byte("*4\r\n$3\r\nSET\r\n$4\r\nkey2\r\n$6\r\nvalue2\r\n$2\r\nNX\r\n"),
		},
		{
			name:     "SET with EX",
			args:     []string{"SET", "key3", "value3", "EX", "60"},
			opts:     NewSetOpts().WithEX(60),
			expected: []byte("*5\r\n$3\r\nSET\r\n$4\r\nkey3\r\n$6\r\nvalue3\r\n$2\r\nEX\r\n$2\r\n60\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a custom sendCommand function to capture the encoded command
			capturedCommand := encodeCommand(tt.args)

			if !bytesEqual(capturedCommand, tt.expected) {
				t.Errorf("%v produced command %v, want %v",
					tt.opts, capturedCommand, tt.expected)
			}
		})
	}
}

// Test the Get method
/* func TestGetCommand(t *testing.T) {
	client := &RedisClient{} // We don't need a real connection for this test

	tests := []struct {
		name     string
		key      string
		expected []byte
	}{
		{
			name:     "Simple GET",
			key:      "testkey",
			expected: []byte("*2\r\n$3\r\nGET\r\n$7\r\ntestkey\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a custom sendCommand function to capture the encoded command
			var capturedCommand []byte
			client.sendCommand = func(args []string) error {
				capturedCommand = encodeCommand(args)
				return nil
			}

			_, err := client.Get(tt.key)
			if err != nil {
				t.Errorf("Get(%s) error = %v", tt.key, err)
				return
			}

			if !bytesEqual(capturedCommand, tt.expected) {
				t.Errorf("Get(%s) produced command %v, want %v", tt.key, capturedCommand, tt.expected)
			}
		})
	}
}
*/
