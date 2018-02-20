package i18n

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/i18n"
)

// Open just walk in input path for preloading localization files
func Open(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".all.yaml") {
			return err
		}

		return i18n.LoadTranslationFile(path)
	})
}
