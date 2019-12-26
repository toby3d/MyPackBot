//go:generate gotext -dir=./../../ -srclang=en update -out=catalog.go -lang=en,ru .
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"gitlab.com/toby3d/mypackbot/internal"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func init() {
	_ = message.Set(language.English, "🕵 Found %d result(s)", plural.Selectf(1, "%d",
		"zero", "🤷 Nothing found",
		"one", "🕵 Found %d result",
		"many", "🕵 Found %[1]d results",
	))

	_ = message.Set(language.Russian, "🕵 Found %d result(s)", plural.Selectf(1, "%d",
		"zero", "🤷 Ничего не найдено",
		"one", "🕵 Найден %[1]d результат",
		"<5", "🕵 Найдено %[1]d результата",
		"many", "🕵 Найдено %[1]d результатов",
	))
}

func main() {
	flagConfig := flag.String("config", filepath.Join(".", "config.yaml"), "set specific path to config")

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
