package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadSettings() {
	log.Debug("Loading .env")
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	log.Debug("Loaded .env")

	env := viper.Get("ENV").(string)

	log.Info(fmt.Sprintf("Utilizando '%s' env", env))

	switch env {
	case "DEVELOPMENT":
		log.SetLevel(log.DebugLevel)
	case "PRODUCTION":
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(log.InfoLevel)
	default:
		log.Fatal("ENV precisa ser ['DEVELOPMENT', 'PRODUCTION']")
	}

	// Non nullable configs
	for _, variable := range []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB", "POSTGRES_PORT"} {
		value := viper.Get(variable)
		if value == nil {
			panic(fmt.Sprintf("Need '%s' variable in .env file", variable))
		}
	}

	// env configs
	for _, variable := range []string{"POSTGRES_HOST", "DEBUG", "max_page", "force_update"} {
		envVariable := fmt.Sprintf("%s_%s", env, variable)
		value := viper.Get(envVariable)
		viper.Set(variable, value)
		if value == nil {
			panic(fmt.Sprintf("Need '%s' variable in .env file", envVariable))
		}
	}

	viper.SetConfigFile("settings.toml")
	err = viper.MergeInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
