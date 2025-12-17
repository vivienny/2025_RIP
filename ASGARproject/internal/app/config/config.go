package config

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string `mapstructure:"service_host"`
	ServicePort int    `mapstructure:"service_port"`
	JWT         JWTConfig
	Redis       RedisConfig
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
	SigningMethod   jwt.SigningMethod // УБРАТЬ *HMAC!
	ExpiresIn       time.Duration
}
type RedisConfig struct {
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	Password    string        `mapstructure:"password"`
	DB          int           `mapstructure:"db"`
	DialTimeout time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout time.Duration `mapstructure:"read_timeout"`
}

func NewConfig() (*Config, error) {
	var err error

	_ = godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("service_host", "0.0.0.0")
	viper.SetDefault("service_port", 8080)
	viper.SetDefault("jwt.secret", "your-super-secret-jwt-key-change-in-production")
	viper.SetDefault("jwt.expiration_hours", 24)
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "redis_password")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.dial_timeout", "10s")
	viper.SetDefault("redis.read_timeout", "10s")

	// Read from environment variables
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("Config file not found, using defaults")
		} else {
			return nil, err
		}
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	// Set JWT values
	cfg.JWT.SigningMethod = jwt.SigningMethodHS256
	cfg.JWT.ExpiresIn = time.Hour * time.Duration(cfg.JWT.ExpirationHours)

	log.Info("config parsed successfully")
	log.Infof("Service: %s:%d", cfg.ServiceHost, cfg.ServicePort)
	log.Infof("JWT Expiration: %d hours", cfg.JWT.ExpirationHours)

	return cfg, nil
}
