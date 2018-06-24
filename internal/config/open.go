package config

import "github.com/spf13/viper"

var Config *viper.Viper

// Open just open configuration file for parsing some data in other functions
func Open(path string) (*viper.Viper, error) {
	cfg := viper.New()

	cfg.AddConfigPath(path)
	cfg.SetConfigName("config")
	cfg.SetConfigType("yaml")

	err := cfg.ReadInConfig()
	return cfg, err
}
