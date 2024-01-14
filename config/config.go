package config

import (
	"github.com/spf13/viper"
	"log"
	"violet/model"
)

var Config model.Config

func ConfigInit() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("[Error] - 找不到配置文件，%s", err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("[Error] - 配置文件读取错误，%s", err)
	}
}
