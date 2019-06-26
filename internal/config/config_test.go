package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDir = filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com", "toby3d", "mypackbot")

func TestOpen(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		cfg, err := Open("/invalid/directory")
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
	t.Run("valid", func(t *testing.T) {
		cfg, err := Open(filepath.Join(testDir, "configs", "config.example.yaml"))
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})
}
