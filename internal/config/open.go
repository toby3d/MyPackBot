package config

import (
	"errors"
	"path/filepath"

	"github.com/spf13/viper"
)

type Reader interface {
	GetString(string) string
	GetInt64(string) int64
}

// Open just open configuration file for parsing some data in other functions
func Open(path string) (*viper.Viper, error) {
	dir, file := filepath.Split(path)
	ext := filepath.Ext(file)

	if file == "" || ext == "" {
		return nil, errors.New("invalid path to config file")
	}

	fileExt := ext[1:]
	fileName := file[:(len(file)-len(fileExt))-1]

	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigName(fileName)
	v.SetConfigType(fileExt)

	err := v.ReadInConfig()

	return v, err
}
