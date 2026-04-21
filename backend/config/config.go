package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppName            string
	AppPort            string
	DBDSN              string
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenMinutes int
	RefreshTokenDays   int
	CookieDomain       string
	CookieSecure       bool
	CookieSameSite     string
}

func LoadConfig() AppConfig {
	_ = godotenv.Load()

	accessMinutes, err := strconv.Atoi(getEnv("ACCESS_TOKEN_MINUTES", "15"))
	if err != nil {
		log.Fatal("invalid ACCESS_TOKEN_MINUTES")
	}

	refreshDays, err := strconv.Atoi(getEnv("REFRESH_TOKEN_DAYS", "7"))
	if err != nil {
		log.Fatal("invalid REFRESH_TOKEN_DAYS")
	}

	cookieSecure, err := strconv.ParseBool(getEnv("COOKIE_SECURE", "false"))
	if err != nil {
		log.Fatal("invalid COOKIE_SECURE")
	}

	cfg := AppConfig{
		AppName:            getEnv("APP_NAME", "go-auth-app"),
		AppPort:            getEnv("APP_PORT", "8080"),
		DBDSN:              os.Getenv("DB_DSN"),
		AccessTokenSecret:  os.Getenv("ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret: os.Getenv("REFRESH_TOKEN_SECRET"),
		AccessTokenMinutes: accessMinutes,
		RefreshTokenDays:   refreshDays,
		CookieDomain:       getEnv("COOKIE_DOMAIN", "localhost"),
		CookieSecure:       cookieSecure,
		CookieSameSite:     getEnv("COOKIE_SAME_SITE", "Lax"),
	}

	if cfg.DBDSN == "" || cfg.AccessTokenSecret == "" || cfg.RefreshTokenSecret == "" {
		log.Fatal("missing required env")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
