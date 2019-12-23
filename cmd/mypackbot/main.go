//go:generate gotext -dir=./../../ -srclang=en update -out=./../../internal/catalog/catalog.go -lang=en,ru .
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"gitlab.com/toby3d/mypackbot/internal"
	_ "gitlab.com/toby3d/mypackbot/internal/catalog"
	"gitlab.com/toby3d/mypackbot/internal/common"
)

func main() {
	flagConfig := flag.String(
		"config", filepath.Join("./", "configs", "config.yaml"), "set specific path to config",
	)

	flag.Parse()

	log.Println("Current build version:", common.Version.String())

	bot, err := internal.New(*flagConfig)
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
	}

	go func() {
		if err = bot.Run(); err != nil {
			log.Fatalln("ERROR:", err.Error())
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
}
