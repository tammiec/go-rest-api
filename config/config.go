package config

import (
	"os"
	"strconv"
	"strings"
)

var (
	Version = "dev"
	Hash    = "dev"
)

type SqlConfig struct {
	Url string
}

type Config struct {
	Name     string
	Version  string
	Hash     string
	HTTPPort int
	HTTPHost string
	DB       SqlConfig
}

func MustGetEnvString(envName string) string {
	value, present := os.LookupEnv(envName)

	if !present || strings.TrimSpace(value) == "" {
		panic("Could not find ENV: " + envName)
	}

	return value
}

func MustGetEnvInt(envName string) int {
	value := MustGetEnvString(envName)

	parsed, err := strconv.Atoi(value)

	if err != nil {
		panic("Could not parse int from ENV: " + envName)
	}

	return parsed
}

func MustGetEnvBool(envName string) bool {
	value := MustGetEnvString(envName)

	parsed, err := strconv.ParseBool(value)

	if err != nil {
		panic("Could not parse bool from ENV: " + envName)
	}

	return parsed
}

func Loader() *Config {
	return &Config{
		Name:     "go-rest-api",
		Version:  Version,
		Hash:     Hash,
		HTTPPort: MustGetEnvInt("HTTP_PORT"),
		HTTPHost: MustGetEnvString("HTTP_HOST"),
		DB: SqlConfig{
			Url: MustGetEnvString("DATABASE_URL"),
		},
	}
}
