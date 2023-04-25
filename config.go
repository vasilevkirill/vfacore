package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func LoadConfig() *viper.Viper {
	Config := viper.New()
	Config.AddConfigPath("./config")
	Config.SetConfigName("config")
	Config.SetConfigType("yaml")

	err := Config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err.Error()))
	}
	return Config
}
