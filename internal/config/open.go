package config

import (
	"path/filepath"

	"github.com/kirillDanshin/dlog"
	"github.com/spf13/viper"
	"gitlab.com/toby3d/mypackbot/internal/errors"
)

type Reader interface {
	GetString(string) string
	GetInt64(string) int64
}

// Open just open configuration file for parsing some data in other functions
func Open(path string) (*viper.Viper, error) {
	dlog.Ln("Opening config on path:", path)

	dir, file := filepath.Split(path)
	ext := filepath.Ext(file)
	if file == "" || ext == "" {
		return nil, errors.New("invalid path to config file")
	}

	fileExt := ext[1:]
	fileName := file[:(len(file)-len(fileExt))-1]

	dlog.Ln("dir:", dir)
	dlog.Ln("file:", file)
	dlog.Ln("fileName:", fileName)
	dlog.Ln("fileExt:", fileExt)

	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigName(fileName)
	v.SetConfigType(fileExt)

	dlog.Ln("Reading", file)
	err := v.ReadInConfig()
	return v, err
}
