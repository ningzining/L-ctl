package config

import (
	"github.com/ningzining/L-ctl/util/templateutil"
	"github.com/spf13/viper"
)

type CtlConfig struct {
	Token string
}

func GetCtlConfig() (*CtlConfig, error) {
	path, err := templateutil.GetConfigFilePath()
	if err != nil {
		return nil, err
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var c CtlConfig
	viper.Unmarshal(&c)
	return &c, nil
}
