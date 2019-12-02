package store

import (
	"sort"
	"strconv"
	"time"

	bolt "github.com/etcd-io/bbolt"
	json "github.com/json-iterator/go"
	"github.com/valyala/fastjson"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"golang.org/x/xerrors"
)

type PhotosStore struct {
	conn     *bolt.DB
	marshler json.API
	parser   fastjson.Parser
}

var (
	ErrPhotoInvalid = model.Error{
		Message: "Invalid photo",
	}

	ErrPhotoExist = model.Error{
		Message: "Photo already exist",
	}

	ErrPhotoNotExist = model.Error{
		Message: "Photo not exist",
	}
)

func NewPhotosStore(conn *bolt.DB, marshler json.API) *PhotosStore {
	return &PhotosStore{
		conn:     conn,
		marshler: marshler,
		parser:   fastjson.Parser{},
	}
}

func (store *PhotosStore) Create(p *model.Photo) error {
	if p == nil || p.FileID == "" {
		return ErrPhotoInvalid
	}

	if store.Get(p.ID) != nil || store.GetByFileID(p.FileID) != nil {
		return ErrPhotoExist
	}

	now := time.Now().UTC().Unix()

	if p.CreatedAt <= 0 {
		p.CreatedAt = now
	}

	if p.UpdatedAt <= 0 {
		p.UpdatedAt = now
	}

	return store.conn.Update(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket(common.BucketPhotos)

		if p.ID, err = bkt.NextSequence(); err != nil {
			return err
		}

		src, err := store.marshler.Marshal(p)
		if err != nil {
			return err
		}

		return bkt.Put([]byte(strconv.FormatUint(p.ID, 10)), src)
	})
}

func (store *PhotosStore) Get(id uint64) *model.Photo {
	p := new(model.Photo)

	if err := store.conn.View(func(tx *bolt.Tx) error {
		return store.marshler.Unmarshal(
			tx.Bucket(common.BucketPhotos).Get([]byte(strconv.FormatUint(id, 10))), p,
		)
	}); err != nil || p.ID == 0 {
		return nil
	}

	return p
}

func (store *PhotosStore) GetByFileID(id string) *model.Photo {
	p := new(model.Photo)

	if err := store.conn.View(func(tx *bolt.Tx) error {
		if err := tx.Bucket(common.BucketPhotos).ForEach(func(key, val []byte) error {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if string(v.GetStringBytes("file_id")) != id {
				return nil
			}

			if err = store.marshler.Unmarshal(val, p); err != nil {
				return err
			}

			return ErrForEachStop
		}); err != nil && !xerrors.Is(err, ErrForEachStop) {
			return err
		}

		return nil
	}); err != nil || p.ID == 0 {
		return nil
	}

	return p
}

func (store *PhotosStore) GetList(offset, limit int) (model.Photos, int) {
	if limit <= 0 {
		limit = 0
	}

	count := 0
	photos := make(model.Photos, 0, limit)
	_ = store.conn.View(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketPhotos).ForEach(func(key, val []byte) (err error) {
			count++

			if count <= offset || (limit > 0 && count > offset+limit) {
				return nil
			}

			p := new(model.Photo)
			if err = store.marshler.Unmarshal(val, p); err != nil {
				return err
			}

			photos = append(photos, p)

			return nil
		})
	})

	sort.Slice(photos, func(i, j int) bool {
		return photos[i].UpdatedAt < photos[j].UpdatedAt
	})

	return photos, count
}

func (store *PhotosStore) Update(p *model.Photo) error {
	if p == nil || p.FileID == "" {
		return ErrPhotoInvalid
	}

	if store.Get(p.ID) == nil && store.GetByFileID(p.FileID) == nil {
		return store.Create(p)
	}

	if p.UpdatedAt <= 0 {
		p.UpdatedAt = time.Now().UTC().Unix()
	}

	src, err := store.marshler.Marshal(p)
	if err != nil {
		return err
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketPhotos).Put([]byte(strconv.FormatUint(p.ID, 10)), src)
	})
}

func (store *PhotosStore) Remove(id uint64) error {
	if store.Get(id) == nil {
		return ErrPhotoNotExist
	}

	return store.conn.Update(func(tx *bolt.Tx) (err error) {
		if err = tx.Bucket(common.BucketPhotos).Delete([]byte(strconv.FormatUint(id, 10))); err != nil {
			return err
		}

		bkt := tx.Bucket(common.BucketUsersPhotos)
		if err = bkt.ForEach(func(key, val []byte) (err error) {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if v.GetUint64("photo_id") != id {
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

func (store *PhotosStore) GetOrCreate(p *model.Photo) (photo *model.Photo, err error) {
	if photo = store.GetByFileID(p.FileID); photo != nil {
		return photo, nil
	}

	if photo = store.Get(p.ID); photo != nil {
		return photo, nil
	}

	if err = store.Create(p); err != nil {
		return nil, err
	}

	return store.GetOrCreate(p)
}
