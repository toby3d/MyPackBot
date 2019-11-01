package main

import (
	"flag"
	"log"
	"path/filepath"

	bunt "github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/migrator"
	"gitlab.com/toby3d/mypackbot/internal/store"
	tg "gitlab.com/toby3d/telegram"
)

var (
	flagOld   = flag.String("old", filepath.Join(".", "old.db"), "filepath to old database file")
	flagNew   = flag.String("new", filepath.Join(".", "new.db"), "filepath to new database file")
	flagToken = flag.String("token", "", "bot token")
)

func main() {
	flag.Parse()

	bot, err := tg.New(*flagToken)
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
	}

	oldDB, err := bunt.Open(*flagOld)
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
	}
	defer oldDB.Close()

	newDB, err := db.Open(*flagNew)
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
	}
	defer newDB.Close()

	if err = migrator.AutoMigrate(migrator.AutoMigrateConfig{
		OldDB: oldDB,
		NewDB: store.NewStore(newDB),
		Bot:   bot,
	}); err != nil {
		log.Fatalln("ERROR:", err.Error())
	}

	log.Println("Done!")
}
