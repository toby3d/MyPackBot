//go:generate gotext -srclang=en update -out=i18n_gen.go -lang=en
package main

import (
	"flag"
	"path/filepath"

	"gitlab.com/toby3d/mypackbot/internal"
)

var flagConfig = flag.String(
	"config",
	filepath.Join("./", "configs", "config.yaml"),
	"set specific path to config",
)

func main() {
	flag.Parse()

	bot, err := internal.New(*flagConfig)
	if err != nil {
		panic(err.Error())
	}

	_ = bot
}
