package config

import (
	"net"
	"strings"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
type Config struct {
	ServerPort         string   `mapstructure:"PORT"`
	MongoURI           string   `mapstructure:"MONGO_URI"`
	DBName             string   `mapstructure:"DB_NAME"`
	JWTSecretKey       string   `mapstructure:"JWT_SECRET_KEY"`
	JWTExpirationHours int      `mapstructure:"JWT_EXPIRATION_HOURS"`
	EnableCache        bool     `mapstructure:"ENABLE_CACHE"`
	RedisAddr          string   `mapstructure:"REDIS_ADDR"`
	RedisHost          string   `mapstructure:"REDIS_HOST"`
	RedisPort          string   `mapstructure:"REDIS_PORT"`
	RedisPassword      string   `mapstructure:"REDIS_PASSWORD"`
	LogLevel           string   `mapstructure:"LOG_LEVEL"`
	LogFormat          string   `mapstructure:"LOG_FORMAT"`
	CookieDomains      []string `mapstructure:"COOKIE_DOMAINS"`
	SecureCookie       bool     `mapstructure:"SECURE_COOKIE"`
	AllowedOrigins     []string `mapstructure:"ALLOWED_ORIGINS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENABLE_CACHE", false)
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("JWT_EXPIRATION_HOURS", 72)
	viper.SetDefault("COOKIE_DOMAINS", []string{"localhost"})
	viper.SetDefault("SECURE_COOKIE", false)
	viper.SetDefault("ALLOWED_ORIGINS", []string{"http://localhost:5173"})

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	// Manually handle comma-separated strings for slices if viper didn't split them
	if allowedOrigins := viper.GetString("ALLOWED_ORIGINS"); allowedOrigins != "" {
		parts := strings.Split(allowedOrigins, ",")
		var cleaned []string
		for _, p := range parts {
			// Trim spaces and quotes
			trimmed := strings.TrimSpace(p)
			trimmed = strings.Trim(trimmed, "\"'")
			if trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}
		config.AllowedOrigins = cleaned
	}

	if cookieDomains := viper.GetString("COOKIE_DOMAINS"); cookieDomains != "" {
		parts := strings.Split(cookieDomains, ",")
		var cleaned []string
		for _, p := range parts {
			// Trim spaces and quotes
			trimmed := strings.TrimSpace(p)
			trimmed = strings.Trim(trimmed, "\"'")
			if trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}
		config.CookieDomains = cleaned
	}

	return
}

// RedisAddress supports the assessment's REDIS_HOST contract while retaining
// REDIS_ADDR compatibility for local development.
func (config Config) RedisAddress() string {
	if config.RedisAddr != "" {
		return config.RedisAddr
	}
	return net.JoinHostPort(config.RedisHost, config.RedisPort)
}
