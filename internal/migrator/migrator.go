package migrator

import (
	"strconv"
	"strings"
	"time"

	"github.com/kirillDanshin/dlog"
	bunt "github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	tg "gitlab.com/toby3d/telegram"
)

type AutoMigrateConfig struct {
	FromPath string
	ToConn   store.Manager
	Bot      *tg.Bot
}

const (
	partState       string = "state"
	partSet         string = "set"
	partSticker     string = "sticker"
	uploadedSetName string = "?"
)

func AutoMigrate(cfg AutoMigrateConfig) error {
	// NOTE(toby3d): open old buntdb database
	db, err := bunt.Open(cfg.FromPath)
	if err != nil {
		dlog.Ln("ERROR:", err.Error())
		return err
	}

	// NOTE(toby3d): preparing temp-stores for migrating
	users := make(map[int]*model.User)   // NOTE(toby3d): users[userId]*User
	sets := make(map[string]struct{})    // NOTE(toby3d): sets[setName]
	userSets := make(map[int][]string)   // NOTE(toby3d): userSets[userId][]setName
	userStickers := make(map[int]string) // NOTE(toby3d): userStickers[userId]fileId

	// NOTE(toby3d): read every key in buntdb database
	if err = db.View(func(tx *bunt.Tx) error {
		return tx.AscendKeys("*", func(k, v string) bool {
			// NOTE(toby3d): split key name on parts
			parts := strings.Split(k, ":")
			// NOTE(toby3d): this part always contains user/chat id
			uid, err := strconv.Atoi(parts[1])
			if err != nil || uid == 0 {
				return true
			}

			// NOTE(toby3d): get now timestamp
			timeStamp := time.Now().UTC().Unix()

			// NOTE(toby3d): we don't modify and force save this data to a new store because keys may be
			// duplicated

			if _, ok := users[uid]; !ok {
				users[uid] = &model.User{
					ID:        uid,
					CreatedAt: timeStamp,
					UpdatedAt: timeStamp,

					LanguageCode: "en",
					LastSeen:     timeStamp,
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
	if err = db.Close(); err != nil {
		return err
	}

	// NOTE(toby3d): STEP 1: migrate users
	for _, u := range users {
		if _, err = cfg.ToConn.Users().GetOrCreate(u); err != nil {
			continue
		}
		dlog.Ln("MIGRATOR: user", u.ID, "successfuly migrated")
	}

	// NOTE(toby3d): STEP 2: migrate sets
	for setName, _ := range sets {
		set, err := cfg.Bot.GetStickerSet(setName)
		if err != nil {
			continue
		}

		for _, sticker := range set.Stickers {
			if _, err = cfg.ToConn.Stickers().GetOrCreate(&model.Sticker{
				ID:         sticker.FileID,
				CreatedAt:  time.Now().UTC().Unix(),
				Width:      sticker.Width,
				Height:     sticker.Height,
				IsAnimated: sticker.IsAnimated,
				SetName:    setName,
				Emoji:      sticker.Emoji,
			}); err != nil {
				continue
			}
			dlog.Ln("MIGRATOR: sticker", sticker.FileID, "successfuly migrated")
		}
	}

	// NOTE(toby3d): STEP 3: import sets to users
	for uid, sets := range userSets {
		for _, setName := range sets {
			u, err := cfg.ToConn.Users().GetOrCreate(users[uid])
			if err != nil || u == nil {
				continue
			}
			if err = cfg.ToConn.AddStickersSet(u, setName); err != nil {
				continue
			}
			dlog.Ln("MIGRATOR: set", setName, "successfuly assigned to", u.ID)
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
			dlog.Ln("MIGRATOR: sticker", fileID, "successfuly sended to", uid)
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
