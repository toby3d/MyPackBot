package migrator

import (
	"strconv"
	"strings"
	"time"

	bunt "github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	tg "gitlab.com/toby3d/telegram"
)

type AutoMigrateConfig struct {
	OldDB *bunt.DB
	NewDB store.Manager
	Bot   *tg.Bot
}

const (
	partState       string = "state"
	partSet         string = "set"
	partSticker     string = "sticker"
	uploadedSetName string = "?"
)

func AutoMigrate(cfg AutoMigrateConfig) (err error) {
	// NOTE(toby3d): preparing temp-stores for migrating
	users := make(map[int]*model.User)   // NOTE(toby3d): users[userId]*User
	sets := make(map[string]struct{})    // NOTE(toby3d): sets[setName]
	userSets := make(map[int][]string)   // NOTE(toby3d): userSets[userId][]setName
	userStickers := make(map[int]string) // NOTE(toby3d): userStickers[userId]fileId

	// NOTE(toby3d): read every key in buntdb database
	if err = cfg.OldDB.View(func(tx *bunt.Tx) error {
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

			if _, ok := users[uid]; !ok {
				users[uid] = new(model.User)
				users[uid] = &model.User{
					ID:           uid,
					LanguageCode: "en",
					LastSeen:     time.Now().UTC().Unix(),
				}
			}

			if strings.EqualFold(parts[2], partSet) {
				setName := parts[3]
				if strings.EqualFold(setName, uploadedSetName) {
					userStickers[uid] = parts[5]
					return true
				}

				if _, ok := sets[setName]; !ok {
					sets[setName] = struct{}{}
				}

				if !contains(userSets[uid], setName) {
					userSets[uid] = append(userSets[uid], setName)
				}
			}

			return true
		})
	}); err != nil {
		return err
	}

	// NOTE(toby3d): close old database
	if err = cfg.OldDB.Close(); err != nil {
		return err
	}

	// NOTE(toby3d): STEP 1: migrate users
	for _, u := range users {
		if _, err = cfg.NewDB.Users().GetOrCreate(u); err != nil {
			continue
		}
	}

	// NOTE(toby3d): STEP 2: migrate sets
	for setName, _ := range sets {
		set, err := cfg.Bot.GetStickerSet(setName)
		if err != nil {
			continue
		}

		for _, sticker := range set.Stickers {
			if _, err = cfg.NewDB.Stickers().GetOrCreate(&model.Sticker{
				ID:         sticker.FileID,
				Width:      sticker.Width,
				Height:     sticker.Height,
				IsAnimated: sticker.IsAnimated,
				SetName:    setName,
				Emoji:      sticker.Emoji,
			}); err != nil {
				continue
			}
		}
	}

	// NOTE(toby3d): STEP 3: import sets to users
	for uid, sets := range userSets {
		for _, setName := range sets {
			u, err := cfg.NewDB.Users().GetOrCreate(users[uid])
			if err != nil || u == nil {
				continue
			}
			if err = cfg.NewDB.AddStickersSet(u, setName); err != nil {
				continue
			}
		}
	}

	// NOTE(toby3d): STEP 4: send uploaded stickers directly to users
	ticker := time.NewTicker(500 * time.Millisecond)
	var count int
	for uid, fileID := range userStickers {
		select {
		case <-ticker.C:
			count++
			if _, err = cfg.Bot.SendSticker(&tg.SendStickerParameters{
				ChatID:              int64(uid),
				Sticker:             fileID,
				DisableNotification: true,
			}); err != nil {
				continue
			}
			if count == len(userStickers) {
				ticker.Stop()
			}
		}
	}

	return nil
}

// contains checks what src array contains find string (or not)
func contains(src []string, find string) bool {
	for i := range src {
		if !strings.EqualFold(src[i], find) {
			continue
		}
		return true
	}
	return false
}
