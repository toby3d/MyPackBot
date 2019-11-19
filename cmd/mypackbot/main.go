//go:generate gotext -srclang=en update -out=./i18n_gen.go -lang=en,ru .
//nolint: gochecknoglobals
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"syscall"

	"gitlab.com/toby3d/mypackbot/internal"
	"gitlab.com/toby3d/mypackbot/internal/common"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flagConfig = flag.String(
		"config", filepath.Join("./", "configs", "config.yaml"), "set specific path to config",
	)
)

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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
