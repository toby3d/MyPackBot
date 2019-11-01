//go:generate gotext -srclang=en update -out=./i18n_gen.go -lang=en,ru .
//nolint: gochecknoglobals
package main

import (
	"flag"
	"log"
	"path/filepath"

	"gitlab.com/toby3d/mypackbot/internal"
)

var (
	gitCommit string

	flagConfig = flag.String(
		"config", filepath.Join("./", "configs", "config.yaml"), "set specific path to config",
	)
	flagV = flag.Bool("v", false, "print current version of the build")
)

func main() {
	flag.Parse()

	if flagV != nil {
		log.Println("Current build version:", gitCommit)
	}

	bot, err := internal.New(*flagConfig)
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
	}

	if err = bot.Run(); err != nil {
		log.Fatalln("ERROR:", err.Error())
	}
}
