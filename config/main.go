package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Signer     map[string]interface{}
	Chain      string
	PayoutRate float64
}

func NewConfig() *viper.Viper {
	cfg := viper.New()
	cfg.SetConfigName("config")
	cfg.SetConfigType("yaml")
	cfg.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file %w", err))
	}

	return cfg
}

func GetConfig(cfg *viper.Viper) Config {
	return Config{
		Signer:     cfg.GetStringMap("Signer"),
		Chain:      cfg.GetString("Chain"),
		PayoutRate: cfg.GetFloat64("PayoutRate"),
	}
}
