package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server server
}

type server struct {
	Address string
	Port    string
}

func GetConfig(filepath string) *Config {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(filepath) // path to look for the config file in
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	var c Config
	err = viper.Unmarshal(&c)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
	return &c
}
