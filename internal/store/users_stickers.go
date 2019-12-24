package store

import (
	"sort"
	"strconv"
	"strings"
	"time"

	bolt "github.com/etcd-io/bbolt"
	json "github.com/json-iterator/go"
	"github.com/valyala/fastjson"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/stickers"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
	"golang.org/x/xerrors"
)

type UsersStickersStore struct {
	conn     *bolt.DB
	marshler json.API
	parser   fastjson.Parser
	stickers stickers.Manager
	users    users.Manager
}

var (
	ErrUserStickerExist = model.Error{
		Message: "Sticker already imported",
	}

	ErrUserStickerNotExist = model.Error{
		Message: "Sticker already removed",
	}
)

func NewUsersStickersStore(db *bolt.DB, us users.Manager, ss stickers.Manager, m json.API) *UsersStickersStore {
	return &UsersStickersStore{
		conn:     db,
		marshler: m,
		parser:   fastjson.Parser{},
		stickers: ss,
		users:    us,
	}
}

func (store *UsersStickersStore) Add(us *model.UserSticker) (err error) {
	if us == nil || us.UserID == 0 || us.StickerID == 0 {
		return nil
	}

	userSticker := store.Get(us)
	if userSticker != nil {
		return ErrUserStickerExist
	}

	timeStamp := time.Now().UTC().Unix()

	if us.CreatedAt == 0 {
		us.CreatedAt = timeStamp
	}

	if us.UpdatedAt == 0 {
		us.UpdatedAt = timeStamp
	}

	return store.conn.Update(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket(common.BucketUsersStickers)

		us.ID, err = bkt.NextSequence()
		if err != nil {
			return err
		}

		src, err := store.marshler.Marshal(us)
		if err != nil {
			return err
		}

		return bkt.Put([]byte(strconv.FormatUint(us.ID, 10)), src)
	})
}

func (store *UsersStickersStore) AddSet(uid uint64, setName string) (err error) {
	set, _ := store.stickers.GetSet(setName)
	for i := range set {
		_ = store.Add(&model.UserSticker{
			UserID:    uid,
			StickerID: set[i].ID,
		})
	}

	return err
}

func (store *UsersStickersStore) Update(us *model.UserSticker) (err error) {
	if us == nil || us.UserID == 0 || us.StickerID == 0 {
		return nil
	}

	userSticker := store.Get(us)
	if userSticker == nil {
		return store.Add(us)
	}

	if us.ID == 0 {
		us.ID = userSticker.ID
	}

	if us.CreatedAt == 0 {
		us.CreatedAt = userSticker.CreatedAt
	}

	if us.UpdatedAt <= userSticker.UpdatedAt {
		us.UpdatedAt = time.Now().UTC().Unix()
	}

	src, err := store.marshler.Marshal(us)
	if err != nil {
		return err
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketUsersStickers).Put([]byte(strconv.FormatUint(us.ID, 10)), src)
	})
}

func (store *UsersStickersStore) Get(us *model.UserSticker) *model.UserSticker {
	if us == nil || us.UserID == 0 || us.StickerID == 0 {
		return nil
	}

	userSticker := new(model.UserSticker)
	if err := store.conn.View(func(tx *bolt.Tx) (err error) {
		if err = tx.Bucket(common.BucketUsersStickers).ForEach(func(key, val []byte) (err error) {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if v.GetUint64("user_id") != us.UserID || v.GetUint64("sticker_id") != us.StickerID {
				return nil
			}

			if err = store.marshler.Unmarshal(val, userSticker); err != nil {
				return err
			}

			return ErrForEachStop
		}); err != nil && !xerrors.Is(err, ErrForEachStop) {
			return err
		}

		return nil
	}); err != nil || userSticker.UserID == 0 || userSticker.StickerID == 0 {
		return nil
	}

	return userSticker
}

//nolint: gocognit
func (store *UsersStickersStore) GetList(uid uint64, offset, limit int, q string) (list model.Stickers, count int) {
	if limit <= 0 {
		limit = 0
	}

	list = make(model.Stickers, 0, limit)
	_ = store.conn.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(common.BucketStickers)
		return tx.Bucket(common.BucketUsersStickers).ForEach(func(key, val []byte) error {
			v, err := store.parser.ParseBytes(val)
			if err != nil || v.GetUint64("user_id") != uid {
				return err
			}

			src := bkt.Get([]byte(strconv.FormatUint(v.GetUint64("sticker_id"), 10)))
			if q != "" {
				vQuery := string(v.GetStringBytes("query"))
				switch {
				case vQuery != "" && !strings.ContainsAny(vQuery, q):
					return nil
				case vQuery == "":
					s, err := store.parser.ParseBytes(src)
					if err != nil || !strings.ContainsAny(string(s.GetStringBytes("emoji")), q) {
						return err
					}
				}
			}

			if (offset != 0 && count <= offset) || (limit > 0 && count > offset+limit) {
				return nil
			}

			count++

			s := new(model.Sticker)
			if err = store.marshler.Unmarshal(src, s); err != nil {
				return err
			}

			list = append(list, s)

			return nil
		})
	})

	sort.Slice(list, func(i, j int) bool {
		return list[i].SetName < list[j].SetName || list[i].UpdatedAt < list[j].UpdatedAt
	})

	return list, count
}

func (store *UsersStickersStore) GetSet(uid uint64, offset, limit int, setName string) (model.Stickers, int) {
	count := 0
	stickers := make(model.Stickers, 0, limit)

	_ = store.conn.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(common.BucketStickers)
		return tx.Bucket(common.BucketUsersStickers).ForEach(func(key, val []byte) (err error) {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if v.GetUint64("user_id") != uid {
				return nil
			}

			src := bkt.Get([]byte(strconv.FormatUint(v.GetUint64("sticker_id"), 10)))
			if v, err = store.parser.ParseBytes(src); err != nil {
				return err
			}

			if !strings.EqualFold(string(v.GetStringBytes("set_name")), setName) {
				return nil
			}

			count++

			if count <= offset || count > limit {
				return nil
			}

			s := new(model.Sticker)
			if err = store.marshler.Unmarshal(src, s); err != nil {
				return err
			}

			stickers = append(stickers, s)

			return nil
		})
	})

	sort.Slice(stickers, func(i, j int) bool {
		return stickers[i].UpdatedAt < stickers[j].UpdatedAt
	})

	return stickers, count
}

func (store *UsersStickersStore) Remove(us *model.UserSticker) (err error) {
	userSticker := store.Get(us)
	if userSticker == nil {
		return ErrUserStickerNotExist
	}

	return store.conn.Update(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket(common.BucketUsersStickers)
		if err = bkt.ForEach(func(key, val []byte) (err error) {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if userSticker.UserID != v.GetUint64("user_id") ||
				(userSticker.StickerID != 0 && userSticker.StickerID != v.GetUint64("sticker_id")) {
				return nil
			}

			if err = bkt.Delete(key); err != nil {
				return err
			}

			return ErrForEachStop
		}); err != nil && !xerrors.Is(err, ErrForEachStop) {
			return err
		}

		return nil
	})
}

func (store *UsersStickersStore) RemoveSet(uid uint64, setName string) (err error) {
	set, _ := store.stickers.GetSet(setName)
	for i := range set {
		_ = store.Remove(&model.UserSticker{
			UserID:    uid,
			StickerID: set[i].ID,
		})
	}

	return err
}
