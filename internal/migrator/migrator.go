package migrator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	json "github.com/json-iterator/go"
	bunt "github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/stickers"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
	usersstickers "gitlab.com/toby3d/mypackbot/internal/model/users/stickers"
	tg "gitlab.com/toby3d/telegram"
)

type (
	Data struct {
		Records Records
		*Backup
	}

	Record struct {
		UserID  int
		SetName string
		FileID  string
		Emoji   string
	}

	Records []*Record

	Backup struct {
		Users        []int    `json:"users"`
		Stickers     []string `json:"stickers"`
		ImportedSets []string `json:"imported_sets"`
		BlockedSets  []string `json:"blocked_sets"`
	}

	AutoMigrateConfig struct {
		OldDB         *bunt.DB
		Stickers      stickers.Manager
		UsersStickers usersstickers.ReadWriter
		Users         users.Manager
		Bot           *tg.Bot
		GroupID       int64
		Marshler      json.API
	}
)

const (
	partSet         string = "set"
	partSticker     string = "sticker"
	uploadedSetName string = "?"
)

func AutoMigrate(cfg AutoMigrateConfig) (err error) {
	// NOTE(toby3d): preparing temp-stores for migrating
	data, err := cfg.importOldData()
	if err != nil {
		return err
	}

	// NOTE(toby3d): STEP 1: migrate users
	if err = cfg.migrateUsers(data); err != nil {
		return err
	}

	// NOTE(toby3d): STEP 2: migrate sets
	if err = cfg.migrateSets(data); err != nil {
		return err
	}

	// NOTE(toby3d): STEP 3: migrate stickers
	return cfg.migrateStickers(data)
}

func (cfg *AutoMigrateConfig) importOldData() (data *Data, err error) {
	data = new(Data)

	if data.Backup, err = cfg.readBackup(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if err = cfg.OldDB.View(func(tx *bunt.Tx) error {
		// NOTE(toby3d): read every key in buntdb database
		return tx.AscendKeys("user:*", func(key, val string) bool {
			r := new(Record)

			// NOTE(toby3d): split key name on parts
			parts := strings.Split(key, ":")

			// NOTE(toby3d): this part always contains user/chat id
			var err error
			r.UserID, err = strconv.Atoi(parts[1])
			if err != nil || r.UserID == 0 || !strings.EqualFold(parts[2], partSet) {
				return true
			}

			switch parts[2] {
			case partSet:
				r.SetName = parts[3]
				r.FileID = parts[5]
			case partSticker:
				r.SetName = common.SetNameUploaded
				r.FileID = parts[3]
			default:
				return true
			}

			if containsString(data.ImportedSets, r.SetName) || containsString(data.Stickers, r.FileID) {
				return true
			}

			if containsString(data.BlockedSets, r.SetName) || r.SetName == uploadedSetName {
				r.SetName = common.SetNameUploaded
			}

			r.Emoji = val

			data.Records = append(data.Records, r)
			return true
		})
	}); err != nil {
		return nil, err
	}

	sort.Slice(data.Records, func(i, j int) bool {
		return data.Records[i].UserID < data.Records[j].UserID ||
			data.Records[i].SetName < data.Records[j].SetName ||
			data.Records[i].FileID < data.Records[j].FileID
	})

	if err = cfg.saveBackup(data.Backup); err != nil {
		return nil, err
	}

	return data, nil
}

func (cfg *AutoMigrateConfig) migrateUsers(data *Data) (err error) {
	if data.Backup, err = cfg.readBackup(); err != nil && !os.IsNotExist(err) {
		return err
	}

	for i := range data.Records {
		if containsInt(data.Users, data.Records[i].UserID) {
			continue
		}

		now := time.Now().UTC().Unix()

		_ = cfg.Users.Create(&model.User{
			ID:           data.Records[i].UserID,
			LanguageCode: "en",
			CreatedAt:    now,
			UpdatedAt:    now,
		})

		data.Users = append(data.Users, data.Records[i].UserID)
		_ = cfg.saveBackup(data.Backup)
	}

	return nil
}

func (cfg *AutoMigrateConfig) migrateSets(data *Data) (err error) {
	if data.Backup, err = cfg.readBackup(); err != nil && !os.IsNotExist(err) {
		return err
	}

	for i := range data.Records {
		if data.Records[i].SetName == uploadedSetName || data.Records[i].SetName == common.SetNameUploaded {
			data.Records[i].SetName = common.SetNameUploaded
			continue
		}

		if containsString(data.ImportedSets, data.Records[i].SetName) {
			continue
		}

		if containsString(data.BlockedSets, data.Records[i].SetName) {
			data.Records[i].SetName = common.SetNameUploaded
			continue
		}

		set, err := cfg.Bot.GetStickerSet(data.Records[i].SetName)
		if err != nil {
			data.BlockedSets = append(data.BlockedSets, data.Records[i].SetName)
			data.Records[i].SetName = common.SetNameUploaded
			_ = cfg.saveBackup(data.Backup)

			continue
		}

		for _, setSticker := range set.Stickers {
			setSticker := setSticker
			_ = cfg.Stickers.Create(stickerToModel(setSticker))
		}

		u := cfg.Users.Get(data.Records[i].UserID)
		_ = cfg.UsersStickers.AddSet(u.ID, set.Name)
		data.ImportedSets = append(data.ImportedSets, set.Name)
		_ = cfg.saveBackup(data.Backup)
	}

	return nil
}

func (cfg *AutoMigrateConfig) migrateStickers(data *Data) (err error) {
	if data.Backup, err = cfg.readBackup(); err != nil && !os.IsNotExist(err) {
		return err
	}

	for i := range data.Records {
		if data.Records[i].SetName == uploadedSetName {
			data.Records[i].SetName = common.SetNameUploaded
		}

		if data.Records[i].SetName != common.SetNameUploaded {
			continue
		}

		if containsString(data.Stickers, data.Records[i].FileID) {
			continue
		}

		// NOTE(toby3d): send old sticker ID to get new
		result, err := cfg.Bot.SendSticker(tg.SendSticker{
			ChatID:              cfg.GroupID,
			DisableNotification: true,
			Sticker:             &tg.InputFile{ID: data.Records[i].FileID},
		})
		if err != nil || !result.IsSticker() {
			continue
		}

		s := stickerToModel(result.Sticker)
		s.SetName = common.SetNameUploaded

		if s.Emoji == "" {
			s.Emoji = data.Records[i].Emoji
		}

		// NOTE(toby3d): store old-new stickers
		_ = cfg.Stickers.Create(s)
		u := cfg.Users.Get(data.Records[i].UserID)
		s = cfg.Stickers.Get(s.FileID)
		_ = cfg.UsersStickers.Add(&model.UserSticker{
			UserID:    u.ID,
			StickerID: s.ID,
		})

		data.Stickers = append(data.Stickers, data.Records[i].FileID)
		_ = cfg.saveBackup(data.Backup)
		_, _ = cfg.Bot.DeleteMessage(result.Chat.ID, result.ID)
	}

	return nil
}

func containsInt(src []int, find int) bool {
	for i := range src {
		if src[i] != find {
			continue
		}

		return true
	}

	return false
}

func containsString(src []string, find string) bool {
	for i := range src {
		if src[i] != find {
			continue
		}

		return true
	}

	return false
}

func stickerToModel(s *tg.Sticker) *model.Sticker {
	sticker := new(model.Sticker)
	sticker.FileID = s.FileID
	sticker.Emoji = s.Emoji
	sticker.Width = s.Width
	sticker.Height = s.Height
	sticker.IsAnimated = s.IsAnimated
	sticker.SetName = s.SetName

	if !sticker.InSet() {
		sticker.SetName = common.SetNameUploaded
	}

	return sticker
}

func (cfg *AutoMigrateConfig) readBackup() (bkp *Backup, err error) {
	bkp = new(Backup)
	filePath := filepath.Join(".", "backup.json")

	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err = cfg.Marshler.Unmarshal(src, bkp); err != nil {
		return nil, err
	}

	return bkp, err
}

func (cfg *AutoMigrateConfig) saveBackup(bkt *Backup) (err error) {
	src, err := cfg.Marshler.Marshal(bkt)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(".", "buckup.json"), src, 0644)
}
