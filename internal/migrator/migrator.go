package migrator

import (
	"strconv"
	"strings"
	"time"

	bunt "github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

type (
	AutoMigrateConfig struct {
		OldDB *bunt.DB
		NewDB store.Manager
		Bot   *tg.Bot
	}

	tempData struct {
		users        map[int]*model.User // NOTE(toby3d): users[userId]*User
		sets         map[string]struct{} // NOTE(toby3d): sets[setName]
		userSets     map[int][]string    // NOTE(toby3d): userSets[userId][]setName
		userStickers map[int]string      // NOTE(toby3d): userStickers[userId]fileId
	}
)

const (
	partSet         string = "set"
	uploadedSetName string = "?"
)

func AutoMigrate(cfg AutoMigrateConfig) (err error) {
	// NOTE(toby3d): preparing temp-stores for migrating
	data, err := importOldData(cfg.OldDB)
	if err != nil {
		return err
	}

	for _, u := range data.users { // NOTE(toby3d): STEP 1: migrate users
		if _, err = cfg.NewDB.Users().GetOrCreate(u); err != nil {
			continue
		}
	}

	for setName := range data.sets { // NOTE(toby3d): STEP 2: migrate sets
		set, err := cfg.Bot.GetStickerSet(setName)
		if err != nil {
			continue
		}

		for _, setSticker := range set.Stickers {
			setSticker := setSticker
			_, _ = cfg.NewDB.Stickers().GetOrCreate(utils.ConvertStickerToModel(&setSticker))
		}
	}

	for uid, sets := range data.userSets { // NOTE(toby3d): STEP 3: import sets to users
		for _, setName := range sets {
			u, err := cfg.NewDB.Users().GetOrCreate(data.users[uid])
			if err != nil || u == nil {
				continue
			}

			_ = cfg.NewDB.AddStickersSet(u, setName)
		}
	}

	count := 0
	ticker := time.NewTicker(500 * time.Millisecond)

	for uid, fileID := range data.userStickers { // NOTE(toby3d): STEP 4: send uploaded stickers directly to users
		for {
			<-ticker.C
			count++

			if _, err = cfg.Bot.SendSticker(&tg.SendStickerParameters{
				ChatID:              int64(uid),
				Sticker:             fileID,
				DisableNotification: true,
			}); err != nil {
				continue
			}

			if count == len(data.userStickers) {
				ticker.Stop()
			}
		}
	}

	return nil
}

func importOldData(db *bunt.DB) (*tempData, error) {
	data := new(tempData)
	err := db.View(func(tx *bunt.Tx) error {
		// NOTE(toby3d): read every key in buntdb database
		return tx.AscendKeys("user:*", func(k, v string) bool {
			// NOTE(toby3d): split key name on parts
			parts := strings.Split(k, ":")
			// NOTE(toby3d): this part always contains user/chat id
			uid, err := strconv.Atoi(parts[1])
			if err != nil || uid == 0 {
				return true
			}

			// NOTE(toby3d): we don't modify and force save this data to a new store because keys may be
			// duplicated

			if _, ok := data.users[uid]; !ok {
				data.users[uid] = new(model.User)
				data.users[uid] = &model.User{
					ID:           uid,
					LanguageCode: "en",
					LastSeen:     time.Now().UTC().Unix(),
				}
			}

			if strings.EqualFold(parts[2], partSet) {
				setName := parts[3]
				if strings.EqualFold(setName, uploadedSetName) {
					data.userStickers[uid] = parts[5]
					return true
				}

				if _, ok := data.sets[setName]; !ok {
					data.sets[setName] = struct{}{}
				}

				if !contains(data.userSets[uid], setName) {
					data.userSets[uid] = append(data.userSets[uid], setName)
				}
			}

			return true
		})
	})

	return data, err
}

// contains checks what src array contains find string (or not)
func contains(src []string, find string) bool {
	var ok bool

	for i := range src {
		if !strings.EqualFold(src[i], find) {
			continue
		}

		ok = true

		break
	}

	return ok
}
