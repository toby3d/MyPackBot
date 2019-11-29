package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	rootPath, err := os.Getwd()
	assert.NoError(t, err)

	t.Run("invalid", func(t *testing.T) {
		cfg, err := Open(filepath.Join("/", "invalid", "directory"))
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
	t.Run("valid", func(t *testing.T) {
		cfg, err := Open(filepath.Join(rootPath, "..", "..", "configs", "config.example.yaml"))
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})
}
