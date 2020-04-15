package config

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

var defaultConfig = []byte(`# filename: config.toml
[webserver]
ports = ["9090"]
contexttimeout = "60s" # timeout for contexts

[userservice]
ports = ["50053"]

[postservice]
ports = ["50052"]

[authservice]
ports = ["50051"]
`)

func NewConfig(filepath string) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("toml")
	viper.AddConfigPath(filepath) // path to look for the config file in
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// if the file was not found just use some default values
			fmt.Println("Config file not found. Using defaults")
			viper.ReadConfig(bytes.NewBuffer(defaultConfig))
		} else {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
}

func NewRuntimeConfig(vip *viper.Viper, filepath string) {
	vip.SetConfigName("config") // name of config file (without extension)
	vip.SetConfigType("toml")
	vip.AddConfigPath(filepath) // path to look for the config file in
	err := vip.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// if the file was not found just use some default values
			fmt.Println("Config file not found. Using defaults")
			vip.ReadConfig(bytes.NewBuffer(defaultConfig))
		} else {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
}
