package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	EnvDevelopment = 4
	EnvTesting     = 2
	EnvProduction  = 0
)

var ErrConfig = errors.New("configuration error")

type Config struct {
	// Environment service running mode enum
	// one of EnvProduction, EnvTesting or EnvDevelopment.
	Environment int
	Postgres    struct {
		Host         string
		Port         int
		User         string
		Password     string
		Database     string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Jwt struct {
		AccessTokenSecretKey           string
		RememberMeAccessTokenSecretKey time.Duration
		AccessTokenDuration            time.Duration
	}
	Smtp struct {
		Sender   string
		Password string
	}
	Authorization struct {
		HeaderKey  string
		PayloadKey string
	}
	RestAddr            string
	RedisAddr           string
	SignUpLinkTokenTime time.Duration
}

func New() *Config {
	cfg := &Config{}
	en := os.Getenv("ENVIRONMENT")

	switch en {
	case "development":
		cfg.Environment = EnvDevelopment
	case "testing":
		cfg.Environment = EnvTesting
	default:
		cfg.Environment = EnvProduction
	}

	return cfg
}

func (cfg *Config) Load(path string) error {
	godotenv.Load(path + "/.env") // load .env file if it exists
	var err error

	cfg.Postgres.Host = os.Getenv("PG_HOST")
	cfg.Postgres.User = os.Getenv("PG_USER")
	cfg.Postgres.Password = os.Getenv("PG_PASSWORD")
	cfg.Postgres.Database = os.Getenv("PG_DATABASE")
	cfg.Jwt.AccessTokenSecretKey = os.Getenv("JWT_ACCESS_TOKEN_SECRET_KEY")
	cfg.Authorization.HeaderKey = os.Getenv("AUTHORIZATION_HEADER_KEY")
	cfg.Authorization.PayloadKey = os.Getenv("AUTHORIZATION_PAYLOAD_KEY")
	cfg.Smtp.Sender = os.Getenv("SMTP_SENDER")
	cfg.Smtp.Password = os.Getenv("SMTP_PASSWORD")

	cfg.RedisAddr = os.Getenv("REDIS_ADDR")

	if cfg.Postgres.MaxIdleTime = os.Getenv("PG_IDLE_TIME"); cfg.Postgres.MaxIdleTime == "" {
		cfg.Postgres.MaxIdleTime = "10m"
	}

	if cfg.Postgres.Port, err = strconv.Atoi(os.Getenv("PG_PORT")); err != nil {
		cfg.Postgres.Port = 5432
	}

	if cfg.Postgres.MaxOpenConns, err = strconv.Atoi(os.Getenv("PG_OPEN_CONNS")); err != nil {
		cfg.Postgres.MaxOpenConns = 4
	}

	if cfg.Postgres.MaxIdleConns, err = strconv.Atoi(os.Getenv("PG_IDLE_CONNS")); err != nil {
		cfg.Postgres.MaxIdleConns = 4
	}

	if cfg.SignUpLinkTokenTime, err = time.ParseDuration(os.Getenv("SIGNUP_LINK_TOKEN_DURATION")); err != nil {
		cfg.SignUpLinkTokenTime = 10 * time.Minute
	}

	if cfg.Jwt.AccessTokenDuration, err = time.ParseDuration(os.Getenv("JWT_ACCESS_TOKEN_DURATION")); err != nil {
		cfg.Jwt.AccessTokenDuration = 12 * time.Hour
	}

	if cfg.Jwt.RememberMeAccessTokenSecretKey, err = time.ParseDuration(os.Getenv("JWT_REMEMBER_ME_ACCESS_TOKEN_DURATION")); err != nil {
		cfg.Jwt.RememberMeAccessTokenSecretKey = time.Hour * 100
	}

	if cfg.RestAddr = os.Getenv("ADDR"); cfg.RestAddr == "" {
		cfg.RestAddr = ":8080"
	}

	return cfg.Validate()
}

func (cfg *Config) Validate() error {
	if cfg.Postgres.Host == "" {
		return fmt.Errorf("%w: postgres host is not set", ErrConfig)
	}

	if cfg.Postgres.User == "" || cfg.Postgres.Password == "" {
		return fmt.Errorf("%w: postgrs username/password is not set", ErrConfig)
	}

	if cfg.Postgres.Database == "" {
		return fmt.Errorf("%w: atabase name is not set", ErrConfig)
	}

	if cfg.RedisAddr == "" {
		return fmt.Errorf("%w: redis address is not set", ErrConfig)
	}

	if cfg.Jwt.AccessTokenSecretKey == "" {
		return fmt.Errorf("%w: jwt secret key is not set", ErrConfig)
	}

	if cfg.Authorization.HeaderKey == "" {
		return fmt.Errorf("%w: header key is not set", ErrConfig)
	}

	if cfg.Authorization.PayloadKey == "" {
		return fmt.Errorf("%w: payload key is not set", ErrConfig)
	}

	if cfg.Smtp.Password == "" {
		return fmt.Errorf("%w: smtp sender password is not set", ErrConfig)
	}

	if cfg.Smtp.Sender == "" {
		return fmt.Errorf("%w: smtp sender is not set", ErrConfig)
	}

	return nil
}

//nolint:nosprintfhostport
func (cfg *Config) PgURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)
}
