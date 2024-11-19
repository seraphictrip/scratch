package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTExpirationInSeconds int
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	if err := godotenv.Load("./.env"); err != nil {
		fmt.Fprintln(os.Stderr, "missing dotenv file, ensure env variables set")
	}

	return Config{
		PublicHost:             EnvOrDefault("HOST", ""),
		Port:                   EnvOrDefault("PORT", "8080"),
		DBUser:                 EnvOrDefault("DB_USER", "root"),
		DBPassword:             EnvOrDefault("DB_PASSWORD", "password"),
		DBAddress:              EnvOrDefault("DB_ADDRESS", "localhost"),
		DBName:                 EnvOrDefault("DB_NAME", "ecom"),
		JWTExpirationInSeconds: EnvIntOrDefault("JWT_EXPIRATION", 3600*24*7),
		JWTSecret:              EnvOrDefault("JWT_SECRET", "not-so-secret"),
	}
}

func MustHaveEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal("missing ", key)
	}
	return val
}

func EnvOrDefault(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func EnvIntOrDefault(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if v, err := strconv.Atoi(val); err == nil {
			return v
		}
	}
	return fallback
}
