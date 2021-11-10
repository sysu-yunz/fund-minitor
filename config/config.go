package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func EnvVariable(key string) string {

	// if .env file exist read from .env
	if _, err := os.Stat("./config/.env"); err == nil {
		return viperEnvVariable(key)
	} else {
		// else read from system
		return os.Getenv(key)
	}
}

// use viper package to read .env file
// return the value of the key
func viperEnvVariable(key string) string {

	//dir, _ := os.Getwd()
	//log.Println(dir)

	// .env - It will search for the .env file in the current directory
	viper.SetConfigFile("./config/.env")

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}
