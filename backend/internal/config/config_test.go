package config

import "testing"

func TestLoadConfigReadsEnvironment(t *testing.T) {
	t.Setenv("MONGO_URI", "mongodb+srv://database.example/test")
	t.Setenv("JWT_SECRET_KEY", "test-signing-key")
	t.Setenv("REDIS_HOST", "redis.example")

	config, err := LoadConfig(t.TempDir())
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.MongoURI != "mongodb+srv://database.example/test" {
		t.Fatalf("MongoURI = %q, expected environment value", config.MongoURI)
	}
	if config.JWTSecretKey != "test-signing-key" {
		t.Fatalf("JWTSecretKey = %q, expected environment value", config.JWTSecretKey)
	}
	if config.RedisHost != "redis.example" {
		t.Fatalf("RedisHost = %q, expected environment value", config.RedisHost)
	}
}

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
