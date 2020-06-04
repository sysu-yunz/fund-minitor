package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

// use viper package to read .env file
// return the value of the key
func ViperEnvVariable(key string) string {

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

func GetWatches() []string {

	return strings.Split("", "-")
}
