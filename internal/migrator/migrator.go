package migrator

import (
	"context"
	"strconv"
	"strings"
	"time"

	bunt "github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/handler"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
	ticker := time.NewTicker(100 * time.Millisecond)
	h := handler.NewHandler(cfg.Bot, cfg.NewDB)

	message.SetString(language.English, "sticker__text", "🤔 This custom/uploaded sticker has been imported from previous version of the bot. You can add it to your pack by clicking on the button below. If the button does not work - please try to click it later when the migration process is completed.")
	message.SetString(language.Russian, "sticker__text", "🤔 Этот загруженый стикер был импортирован с прошлой версии бота. Ты можешь добавить его к себе нажав на кнопку ниже. Если кнопка не работает - пожалуйста, попробуй нажать её позже, когда процесс миграции завершится.")
	message.SetString(language.English, "sticker__button_add-single", "📙 Import this sticker")
	message.SetString(language.Russian, "sticker__button_add-single", "📙 Импортировать этот стикер")
	matcher := language.NewMatcher([]language.Tag{language.English, language.Russian})

	for uid, fileID := range data.userStickers { // NOTE(toby3d): STEP 4: send uploaded stickers directly to users
		count++

		if count > len(data.userStickers) {
			ticker.Stop()
			break
		}

		<-ticker.C

		reply, err := cfg.Bot.SendSticker(&tg.SendStickerParameters{
			ChatID:              int64(uid),
			Sticker:             fileID,
			DisableNotification: true,
		})
		if err != nil {
			continue
		}

		u := cfg.NewDB.Users().Get(uid)
		s := utils.ConvertStickerToModel(reply.Sticker)
		s.CreatedAt = reply.Date
		s.UpdatedAt = reply.Date

		ctx := context.Background()
		tag, _, _ := matcher.Match(language.Make(u.LanguageCode))
		printer := message.NewPrinter(tag)
		ctx = context.WithValue(ctx, common.ContextPrinter, printer)
		ctx = context.WithValue(ctx, common.ContextUser, u)
		ctx = context.WithValue(ctx, common.ContextSticker, s)

		if err = h.IsSticker(ctx, reply); err != nil {
			continue
		}
	}

	return nil
}

func importOldData(db *bunt.DB) (*tempData, error) {
	data := new(tempData)
	data.users = make(map[int]*model.User)
	data.sets = make(map[string]struct{})
	data.userSets = make(map[int][]string)
	data.userStickers = make(map[int]string)

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
