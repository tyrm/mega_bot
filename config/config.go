package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DBEngine string

	DiscordToken string

	LoggerConfig string

	ResponderWorkers int
}

func CollectConfig() *Config {
	var missingEnv []string
	var config Config

	// DB_ENGINE
	config.DBEngine = os.Getenv("DB_ENGINE")
	if config.DBEngine == "" {
		missingEnv = append(missingEnv, "DB_ENGINE")
	}

	// DISCORD_TOKEN
	config.DiscordToken = os.Getenv("DISCORD_TOKEN")

	// LOG_LEVEL
	var envLoggerLevel = os.Getenv("LOG_LEVEL")
	if envLoggerLevel == "" {
		config.LoggerConfig = "<root>=INFO"
	} else {
		config.LoggerConfig = fmt.Sprintf("<root>=%s", strings.ToUpper(envLoggerLevel))
	}

	// LOG_LEVEL
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
