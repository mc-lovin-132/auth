package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	Host string
	Port string

	SecretKey          string
	Salt               string
	AccessTTLMinutes   int
	RefreshTTLDays     int
	StartDeviceVersion int
	StartGlobalVersion int

	UserHost string
	UserPort int

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisBDNumber int
}

func ErrConfigLoadingFailed(err error) error {
	return fmt.Errorf("config loading failed: %w", err)
}

func New() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, ErrConfigLoadingFailed(err)
	}

	cfg := &Config{}

	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, ErrConfigLoadingFailed(err)
	}
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.DBSSLMode = os.Getenv("DB_SSLMODE")

	cfg.Host = os.Getenv("HOST")
	cfg.Port = os.Getenv("PORT")

	cfg.SecretKey = os.Getenv("SECRET_KEY")
	cfg.AccessTTLMinutes, err = strconv.Atoi(os.Getenv("ACCESS_TTL_MINUTES"))
	if err != nil {
		return nil, ErrConfigLoadingFailed(err)
	}
	cfg.RefreshTTLDays, err = strconv.Atoi(os.Getenv("REFRESH_TTL_DAYS"))
	if err != nil {
		return nil, ErrConfigLoadingFailed(err)
	}

	cfg.UserHost = os.Getenv("USER_HOST")
	cfg.UserPort, err = strconv.Atoi(os.Getenv("USER_PORT"))
	if err != nil {
		return nil, ErrConfigLoadingFailed(err)
	}

	cfg.RedisHost = os.Getenv("REDIS_HOST")
	cfg.RedisPort = os.Getenv("REDIS_PORT")
	cfg.RedisPassword = os.Getenv("REDIS_PASSWORD")
	cfg.RedisBDNumber, err = strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))
	if err != nil {
		return nil, ErrConfigLoadingFailed(err)
	}
	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

func (c *Config) Addr() string {
	return c.Host + ":" + c.Port
}

func (c *Config) RedisAddr() string {
	return c.RedisHost + ":" + c.RedisPort
}

func (c *Config) UserAddr() string {
	s := fmt.Sprintf("%s:%d", c.UserHost, c.UserPort)
	fmt.Println(s)
	return s
}
