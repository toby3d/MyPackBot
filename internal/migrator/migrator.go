package migrator

import (
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

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
		UserID  int64
		SetName string
		FileID  string
		Emoji   string
	}

	Records []*Record

	Backup struct {
		Users        []int64  `json:"users"`
		Stickers     []string `json:"stickers"`
		ImportedSets []string `json:"imported_sets"`
		BlockedSets  []string `json:"blocked_sets"`
	}

	AutoMigrateConfig struct {
		OldDB         *bunt.DB
		Stickers      stickers.Manager
		UsersStickers usersstickers.Manager
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
	cfg.migrateUsers(data)

	// NOTE(toby3d): STEP 2: migrate sets
	cfg.migrateSets(data)

	// NOTE(toby3d): STEP 3: migrate stickers
	cfg.migrateStickers(data)

	return nil
}

func (cfg *AutoMigrateConfig) importOldData() (*Data, error) {
	data := new(Data)

	var err error
	if data.Backup, err = cfg.readBackup("backup.json"); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if data.Backup == nil {
		data.Backup = new(Backup)
	}

	data.Records = make([]*Record, 0)

	if err = cfg.OldDB.View(func(tx *bunt.Tx) error {
		// NOTE(toby3d): read every key in buntdb database
		return tx.AscendKeys("user:*", func(key, val string) bool {
			r := new(Record)

			// NOTE(toby3d): split key name on parts
			parts := strings.Split(key, ":")

			// NOTE(toby3d): this part always contains user/chat id
			var err error
			r.UserID, err = strconv.ParseInt(parts[1], 10, 64)
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

	if err = cfg.saveBackup("backup.json", data.Backup); err != nil {
		return nil, err
	}

	return data, nil
}

func (cfg *AutoMigrateConfig) migrateUsers(data *Data) (err error) {
	if data.Backup, err = cfg.readBackup("backup.json"); err != nil && !os.IsNotExist(err) {
		return err
	}

	for i := range data.Records {
		if containsInt(data.Users, data.Records[i].UserID) {
			continue
		}

		cfg.Users.Create(&model.User{
			UserID:       data.Records[i].UserID,
			LanguageCode: "en",
		})

		data.Users = append(data.Users, data.Records[i].UserID)
		cfg.saveBackup("backup.json", data.Backup)
	}

	return nil
}

func (cfg *AutoMigrateConfig) migrateSets(data *Data) (err error) {
	if data.Backup, err = cfg.readBackup("backup.json"); err != nil && !os.IsNotExist(err) {
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
			cfg.saveBackup("backup.json", data.Backup)
			continue
		}

		for _, setSticker := range set.Stickers {
			setSticker := setSticker
			cfg.Stickers.Create(stickerToModel(&setSticker))
		}

		u := cfg.Users.GetByUserID(data.Records[i].UserID)
		cfg.UsersStickers.AddSet(u.ID, set.Name)
		data.ImportedSets = append(data.ImportedSets, set.Name)
		cfg.saveBackup("backup.json", data.Backup)
	}

	return nil
}

func (cfg *AutoMigrateConfig) migrateStickers(data *Data) (err error) {
	if data.Backup, err = cfg.readBackup("backup.json"); err != nil && !os.IsNotExist(err) {
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
		result, err := cfg.Bot.SendSticker(&tg.SendStickerParameters{
			ChatID:              cfg.GroupID,
			DisableNotification: true,
			Sticker:             data.Records[i].FileID,
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
		cfg.Stickers.Create(s)
		u := cfg.Users.GetByUserID(data.Records[i].UserID)
		s = cfg.Stickers.GetByFileID(s.FileID)
		cfg.UsersStickers.Add(u.ID, s.ID)
		data.Stickers = append(data.Stickers, data.Records[i].FileID)
		cfg.saveBackup("backup.json", data.Backup)
		cfg.Bot.DeleteMessage(result.Chat.ID, result.ID)
	}

	return nil
}

func containsInt(src []int64, find int64) bool {
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

func (cfg *AutoMigrateConfig) readBackup(filePath string) (bkp *Backup, err error) {
	bkp = new(Backup)
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

func (cfg *AutoMigrateConfig) saveBackup(filePath string, bkt *Backup) (err error) {
	src, err := cfg.Marshler.Marshal(bkt)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, src, 0644)
}
