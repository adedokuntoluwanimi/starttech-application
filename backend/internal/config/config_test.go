package config

import "testing"

func TestRedisAddress(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name:     "uses explicit legacy address",
			config:   Config{RedisAddr: "redis.example:6380", RedisHost: "ignored", RedisPort: "6379"},
			expected: "redis.example:6380",
		},
		{
			name:     "combines assessment host and port",
			config:   Config{RedisHost: "starttech-redis.example", RedisPort: "6379"},
			expected: "starttech-redis.example:6379",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := test.config.RedisAddress(); actual != test.expected {
				t.Fatalf("RedisAddress() = %q, expected %q", actual, test.expected)
			}
		})
	}
}
