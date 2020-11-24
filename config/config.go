package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	CookieSecret string

	DBEngine string

	DiscordKey    string
	DiscordSecret string
	DiscordToken  string

	ExtHostname string

	LoggerConfig string

	RedisAddress  string
	RedisPassword string

	ResponderWorkers int
}

func CollectConfig() *Config {
	var missingEnv []string
	var config Config

	// COOKIE_SECRET
	config.CookieSecret = os.Getenv("COOKIE_SECRET")
	if config.CookieSecret == "" {
		missingEnv = append(missingEnv, "COOKIE_SECRET")
	}

	// DB_ENGINE
	config.DBEngine = os.Getenv("DB_ENGINE")
	if config.DBEngine == "" {
		missingEnv = append(missingEnv, "DB_ENGINE")
	}

	// DISCORD_KEY
	config.DiscordKey = os.Getenv("DISCORD_KEY")

	// DISCORD_SECRET
	config.DiscordSecret = os.Getenv("DISCORD_SECRET")

	// DISCORD_TOKEN
	config.DiscordToken = os.Getenv("DISCORD_TOKEN")

	// EXT_HOSTNAME
	config.ExtHostname = os.Getenv("EXT_HOSTNAME")
	if config.ExtHostname == "" {
		missingEnv = append(missingEnv, "EXT_HOSTNAME")
	}

	// LOG_LEVEL
	var envLoggerLevel = os.Getenv("LOG_LEVEL")
	if envLoggerLevel == "" {
		config.LoggerConfig = "<root>=INFO"
	} else {
		config.LoggerConfig = fmt.Sprintf("<root>=%s", strings.ToUpper(envLoggerLevel))
	}

	// REDIS_ADDRESS
	config.RedisAddress = os.Getenv("REDIS_ADDRESS")
	if config.RedisAddress == "" {
		missingEnv = append(missingEnv, "REDIS_ADDRESS")
	}

	// REDIS_PASSWORD
	config.RedisPassword = os.Getenv("REDIS_PASSWORD")

	// RESPONDER_WORKERS
	var envResponderWorkers = os.Getenv("RESPONDER_WORKERS")
	if envResponderWorkers == "" {
		config.ResponderWorkers = 4
	} else {
		i, err := strconv.Atoi(envResponderWorkers)
		if err != nil {
			panic(err)
		}
		config.ResponderWorkers = i
	}

	// Validation
	if len(missingEnv) > 0 {
		msg := fmt.Sprintf("Environment variables missing: %v", missingEnv)
		panic(fmt.Sprint(msg))
	}

	return &config
}
