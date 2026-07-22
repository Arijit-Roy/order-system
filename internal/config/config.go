package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisAddr  string
	GRPCAddr   string
}

func Load() Config {
	return Config{
		DBHost:     getenv("DB_HOST", "localhost"),
		DBPort:     getenv("DB_PORT", "5432"),
		DBUser:     getenv("DB_USER", "orderuser"),
		DBPassword: getenv("DB_PASSWORD", "orderpass"),
		DBName:     getenv("DB_NAME", "ordersdb"),

		RedisAddr: getenv("REDIS_ADDR", "localhost:6379"),
		GRPCAddr:  getenv("GRPC_ADDR", "localhost:50051"),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (c Config) PostgresDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)

}
