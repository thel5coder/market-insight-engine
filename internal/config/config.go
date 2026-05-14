package config

import (
	"github.com/spf13/viper"
	"time"
)

type AppConfig struct {
	Env  string
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string

	PoolSize          int
	MinIdleConn       int
	MaxLifetime       time.Duration
	IdleTimeout       time.Duration
	ConnectionTimeout int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string

	MaxActiveConnection int
	MaxIdleConnection   int
	Timeout             time.Duration
	Wait                bool
}

type ExternalApi struct {
	BaseUrl   string
	SecretKey string
	APIKey    string
}

type ExternalMarketAPI struct {
	BinanceApi    ExternalApi
	FinhubApi     ExternalApi
	TwelveDataApi ExternalApi
	TappiApi      ExternalApi
}

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	APIs     ExternalMarketAPI
}

// Load reads the variables and explicitly maps them.
func Load() (*Config, error) {
	vpr := viper.New()
	vpr.SetConfigFile(".env")
	vpr.AutomaticEnv()

	_ = vpr.ReadInConfig()

	vpr.SetDefault("DB_POOL", 10)
	vpr.SetDefault("DB_MIN_IDLE", 2)
	vpr.SetDefault("DB_MAX_LIFETIME", "1h")
	vpr.SetDefault("DB_IDLE_TIMEOUT", "30m")
	vpr.SetDefault("DB_CONNECTION_TIMEOUT", 10)

	vpr.SetDefault("REDIS_MAX_ACTIVE_CONNECTION", 100)
	vpr.SetDefault("REDIS_MAX_IDLE_CONNECTION", 20)
	vpr.SetDefault("REDIS_TIMEOUT", "5s")
	vpr.SetDefault("REDIS_WAIT", true)

	cfg := &Config{
		App: AppConfig{
			Env:  vpr.GetString("APP_ENV"),
			Port: vpr.GetString("APP_PORT"),
		},
		Postgres: PostgresConfig{
			Host:              vpr.GetString("DB_HOST"),
			Port:              vpr.GetString("DB_PORT"),
			User:              vpr.GetString("DB_USER"),
			Password:          vpr.GetString("DB_PASSWORD"),
			Name:              vpr.GetString("DB_NAME"),
			SSLMode:           vpr.GetString("DB_SSLMODE"),
			PoolSize:          vpr.GetInt("DB_POOL"),
			MinIdleConn:       vpr.GetInt("DB_MIN_IDLE"),
			MaxLifetime:       vpr.GetDuration("DB_MAX_LIFETIME"),
			IdleTimeout:       vpr.GetDuration("DB_IDLE_TIMEOUT"),
			ConnectionTimeout: vpr.GetInt("DB_CONNECTION_TIMEOUT"),
		},
		Redis: RedisConfig{
			Host:     vpr.GetString("REDIS_HOST"),
			Port:     vpr.GetString("REDIS_PORT"),
			Password: vpr.GetString("REDIS_PASSWORD"),
		},
		APIs: ExternalMarketAPI{
			BinanceApi: ExternalApi{
				BaseUrl:   vpr.GetString("BINANCE_BASE_URL"),
				SecretKey: vpr.GetString("BINANCE_SECRET_KEY"),
				APIKey:    vpr.GetString("BINANCE_API_KEY"),
			},
		},
	}

	return cfg, nil
}
