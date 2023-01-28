package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"runtime"
)

var configObject ConfigObject

func init() {
	configObject = readViperConfig()
}

func readViperConfig() (obj ConfigObject) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("./params")
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	v.AddConfigPath(fmt.Sprintf("%s/../params", basepath))
	v.SetConfigName("app")
	v.SetConfigType("env")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err == nil {
		log.Printf("Using config file: %s", v.ConfigFileUsed())
	} else {
		log.Panicf("Config error: %s", err)
	}

	err = v.Unmarshal(&obj)
	if err != nil {
		log.Panicf("Config error: %s", err)
	}
	return
}

// Config return provider so that you can read config anywhere
func Config() ConfigObject {
	return configObject
}
