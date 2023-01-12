package utils

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var envVariables = []string{"DEBUG", "max_page", "force_update", "POSTGRES_HOST"}
var postgresVariables = []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB", "POSTGRES_PORT"}

func LoadSettings() {
	if _, err := os.Stat(".env"); err != nil {
		os.Chdir("../")
	}

	log.Debug("Loading .env")
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	log.Debug("Loaded .env")

	env := viper.GetString("ENV")

	log.Info(fmt.Sprintf("Utilizando '%s' env", env))

	format := &log.JSONFormatter{}
	log.SetFormatter(format)

	switch env {
	case "DEVELOPMENT":
		f, err := os.OpenFile(".log", os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(f)
		log.SetLevel(log.DebugLevel)

	case "PRODUCTION":
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(log.InfoLevel)
	default:
		log.Fatal("ENV precisa ser ['DEVELOPMENT', 'PRODUCTION']")
	}

	// Non nullable configs
	for _, variable := range postgresVariables {
		value := viper.Get(variable)
		if value == nil {
			panic(fmt.Sprintf("Need '%s' variable in .env file", variable))
		}
	}

	// env configs
	for _, variable := range envVariables {
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
