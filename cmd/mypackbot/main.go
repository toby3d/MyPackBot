//go:generate gotext -srclang=en update -out=./i18n_gen.go -lang=en,ru .
//nolint: gochecknoglobals
package main

import (
	"flag"
	"log"
	"path/filepath"

	"gitlab.com/toby3d/mypackbot/internal"
	"gitlab.com/toby3d/mypackbot/internal/common"
)

var flagConfig = flag.String(
	"config", filepath.Join("./", "configs", "config.yaml"), "set specific path to config",
)

func main() {
	flag.Parse()

	log.Println("Current build version:", common.Version.String())

	bot, err := internal.New(*flagConfig)
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
	}

	if err = bot.Run(); err != nil {
		log.Fatalln("ERROR:", err.Error())
	}
}
