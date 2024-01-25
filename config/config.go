package config

import (
	"github.com/ningzining/lctl/util/templateutil"
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
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
