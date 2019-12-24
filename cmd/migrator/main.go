package main

import (
	"flag"
	"log"
	"path/filepath"

	json "github.com/json-iterator/go"
	bunt "github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/migrator"
	"gitlab.com/toby3d/mypackbot/internal/store"
	tg "gitlab.com/toby3d/telegram"
)

func main() {
	var (
		flagGroup = flag.Int64("group", 0, "proxy group for migration")
		flagNew   = flag.String("new", filepath.Join(".", "new.db"), "filepath to new database file")
		flagOld   = flag.String("old", filepath.Join(".", "old.db"), "filepath to old database file")
		flagToken = flag.String("token", "", "bot token")
	)

	flag.Parse()

	bot, err := tg.New(*flagToken)
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
	}

	oldDB, err := bunt.Open(*flagOld)
	if err != nil {
		log.Fatalln("ERROR OLD DB:", err.Error())
	}
	defer oldDB.Close()

	newDB, err := db.Open(*flagNew)
	if err != nil {
		log.Fatalln("ERROR NEW DB:", err.Error())
	}
	defer newDB.Close()

	marshler := json.ConfigFastest
	users := store.NewUsersStore(newDB, marshler)
	stickers := store.NewStickersStore(newDB, marshler)
	usersStickers := store.NewUsersStickersStore(newDB, users, stickers, marshler)

	if err = migrator.AutoMigrate(migrator.AutoMigrateConfig{
		Bot:           bot,
		GroupID:       *flagGroup,
		Stickers:      stickers,
		UsersStickers: usersStickers,
		Users:         users,
		OldDB:         oldDB,
		Marshler:      marshler,
	}); err != nil {
		log.Fatalln("ERROR:", err.Error())
	}

	log.Println("Done!")
}
