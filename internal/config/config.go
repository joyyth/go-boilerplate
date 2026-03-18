package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Database DatabaseConfig `koanf:"database" validate:"required"`
	Server   ServerConfig   `koanf:"server"   validate:"required"`
	Logger   LoggerConfig   `koanf:"logger"   validate:"required"`
	Auth     AuthConfig     `koanf:"auth"     validate:"required"`
}
type LoggerConfig struct {
	Level  string `koanf:"level" validate:"required,oneof=debug info warn error"`
	Pretty bool   `koanf:"pretty"`
}
type AuthConfig struct {
	AccessSecret       string        `koanf:"access_secret"        validate:"required"`
	RefreshSecret      string        `koanf:"refresh_secret"       validate:"required"`
	AccessTokenExpiry  time.Duration `koanf:"access_token_expiry"  validate:"required"`
	RefreshTokenExpiry time.Duration `koanf:"refresh_token_expiry" validate:"required"`
}

type DatabaseConfig struct {
	Host            string        `koanf:"host"              validate:"required"`
	Port            int           `koanf:"port"              validate:"required"`
	User            string        `koanf:"user"              validate:"required"`
	Password        string        `koanf:"password"`
	Name            string        `koanf:"name"              validate:"required"`
	SSLMode         string        `koanf:"ssl_mode"          validate:"required,oneof=disable require verify-ca verify-full"`
	MaxOpenConns    int           `koanf:"max_open_conns"    validate:"required,min=1"`
	MaxIdleConns    int           `koanf:"max_idle_conns"    validate:"required,min=1"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime time.Duration `koanf:"conn_max_idle_time" validate:"required"`
}

type ServerConfig struct {
	Port               int      `koanf:"port"                validate:"required,min=1,max=65535"`
	ReadTimeout        int      `koanf:"read_timeout"        validate:"required,min=1"`
	WriteTimeout       int      `koanf:"write_timeout"       validate:"required,min=1"`
	IdleTimeout        int      `koanf:"idle_timeout"        validate:"required,min=1"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins"`
}

func LoadConfig() (*Config, error) {

	k := koanf.New(".")

	err := k.Load(env.Provider(".", env.Opt{
		Prefix: "YOURAPP_",
		TransformFunc: func(k, v string) (string, any) {
			key := strings.TrimPrefix(k, "YOURAPP_")
			key = strings.ToLower(key)
			key = strings.ReplaceAll(key, "__", ".")
			return key, v
		},
	}), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load env config: %w", err)
	}

	cfg := &Config{}
	if err = k.Unmarshal("", cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err = validator.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	return cfg, nil
}
