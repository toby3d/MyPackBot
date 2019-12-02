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
	"golang.org/x/xerrors"
)

type StickersStore struct {
	conn     *bolt.DB
	marshler json.API
	parser   fastjson.Parser
}

var (
	ErrStickerInvalid = model.Error{
		Message: "Invalid sticker",
	}

	ErrStickerExist = model.Error{
		Message: "Sticker already exist",
	}

	ErrStickerNotExist = model.Error{
		Message: "Sticker not exist",
	}
)

func NewStickersStore(conn *bolt.DB, marshler json.API) *StickersStore {
	var parser fastjson.Parser
	return &StickersStore{
		conn:     conn,
		marshler: marshler,
		parser:   parser,
	}
}

func (store *StickersStore) Create(s *model.Sticker) error {
	if s == nil || s.FileID == "" {
		return ErrStickerInvalid
	}

	if store.Get(s.ID) != nil || store.GetByFileID(s.FileID) != nil {
		return ErrStickerExist
	}

	now := time.Now().UTC().Unix()

	if s.CreatedAt <= 0 {
		s.CreatedAt = now
	}

	if s.UpdatedAt <= 0 {
		s.UpdatedAt = now
	}

	if s.SetName == "" {
		s.SetName = common.SetNameUploaded
	}

	return store.conn.Update(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket(common.BucketStickers)

		if s.ID, err = bkt.NextSequence(); err != nil {
			return err
		}

		src, err := store.marshler.Marshal(s)
		if err != nil {
			return err
		}

		return bkt.Put([]byte(strconv.FormatUint(s.ID, 10)), src)
	})
}

func (store *StickersStore) Get(id uint64) *model.Sticker {
	s := new(model.Sticker)

	if err := store.conn.View(func(tx *bolt.Tx) error {
		return store.marshler.Unmarshal(
			tx.Bucket(common.BucketStickers).Get([]byte(strconv.FormatUint(id, 10))), s,
		)
	}); err != nil || s.ID == 0 {
		return nil
	}

	return s
}

func (store *StickersStore) GetByFileID(id string) *model.Sticker {
	s := new(model.Sticker)

	if err := store.conn.View(func(tx *bolt.Tx) error {
		if err := tx.Bucket(common.BucketStickers).ForEach(func(key, val []byte) error {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if string(v.GetStringBytes("file_id")) != id {
				return nil
			}

			if err = store.marshler.Unmarshal(val, s); err != nil {
				return err
			}

			return ErrForEachStop
		}); err != nil && !xerrors.Is(err, ErrForEachStop) {
			return err
		}

		return nil
	}); err != nil || s.ID == 0 {
		return nil
	}

	return s
}

func (store *StickersStore) GetList(offset, limit int, query string) (model.Stickers, int) {
	if limit <= 0 {
		limit = 0
	}

	count := 0
	stickers := make(model.Stickers, 0, limit)
	_ = store.conn.View(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketStickers).ForEach(func(key, val []byte) error {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if query != "" && !strings.ContainsAny(v.Get("emoji").String(), query) {
				return nil
			}

			count++

			if count <= offset || (limit > 0 && count > offset+limit) {
				return nil
			}

			s := new(model.Sticker)
			if err = store.marshler.Unmarshal(val, s); err != nil {
				return err
			}

			stickers = append(stickers, s)

			return nil
		})
	})

	sort.Slice(stickers, func(i, j int) bool {
		return stickers[i].SetName < stickers[j].SetName || stickers[i].UpdatedAt < stickers[j].UpdatedAt
	})

	return stickers, count
}

func (store *StickersStore) GetSet(name string) (model.Stickers, int) {
	count := 0
	stickers := make(model.Stickers, 0)

	_ = store.conn.View(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketStickers).ForEach(func(key, val []byte) (err error) {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if !strings.EqualFold(string(v.GetStringBytes("set_name")), name) {
				return nil
			}

			count++

			s := new(model.Sticker)
			if err = store.marshler.Unmarshal(val, s); err != nil {
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

func (store *StickersStore) Update(s *model.Sticker) error {
	if s == nil || s.FileID == "" {
		return ErrStickerInvalid
	}

	if store.Get(s.ID) == nil && store.GetByFileID(s.FileID) == nil {
		return store.Create(s)
	}

	if s.UpdatedAt <= 0 {
		s.UpdatedAt = time.Now().UTC().Unix()
	}

	if s.SetName == "" {
		s.SetName = common.SetNameUploaded
	}

	src, err := store.marshler.Marshal(s)
	if err != nil {
		return err
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketStickers).Put([]byte(strconv.FormatUint(s.ID, 10)), src)
	})
}

func (store *StickersStore) Remove(id uint64) error {
	if store.Get(id) == nil {
		return ErrStickerNotExist
	}

	return store.conn.Update(func(tx *bolt.Tx) (err error) {
		if err = tx.Bucket(common.BucketStickers).Delete([]byte(strconv.FormatUint(id, 10))); err != nil {
			return err
		}

		bkt := tx.Bucket(common.BucketUsersStickers)

		if err = bkt.ForEach(func(key, val []byte) (err error) {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if v.GetUint64("sticker_id") != id {
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

func (store *StickersStore) GetOrCreate(s *model.Sticker) (sticker *model.Sticker, err error) {
	if sticker = store.GetByFileID(s.FileID); sticker != nil {
		return sticker, nil
	}

	if sticker = store.Get(s.ID); sticker != nil {
		return sticker, nil
	}

	if err = store.Create(s); err != nil {
		return nil, err
	}

	return store.GetOrCreate(s)
}
