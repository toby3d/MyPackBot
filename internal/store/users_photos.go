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
	"gitlab.com/toby3d/mypackbot/internal/model/photos"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
	"golang.org/x/xerrors"
)

type UsersPhotosStore struct {
	conn     *bolt.DB
	marshler json.API
	parser   fastjson.Parser
	photos   photos.Manager
	users    users.Manager
}

var (
	ErrUserPhotoExist = model.Error{
		Message: "Photo already imported",
	}

	ErrUserPhotoNotExist = model.Error{
		Message: "Photo already removed",
	}
)

func NewUsersPhotosStore(conn *bolt.DB, us users.Manager, ps photos.Manager, marshler json.API) *UsersPhotosStore {
	return &UsersPhotosStore{
		conn:     conn,
		marshler: marshler,
		parser:   fastjson.Parser{},
		photos:   ps,
		users:    us,
	}
}

func (store *UsersPhotosStore) Add(up *model.UserPhoto) (err error) {
	if up == nil || up.UserID == 0 || up.PhotoID == 0 {
		return nil
	}

	userPhoto := store.Get(up)
	if userPhoto != nil {
		return ErrUserPhotoExist
	}

	timeStamp := time.Now().UTC().Unix()
	if up.CreatedAt == 0 {
		up.CreatedAt = timeStamp
	}

	if up.UpdatedAt == 0 {
		up.UpdatedAt = timeStamp
	}

	return store.conn.Update(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket(common.BucketUsersPhotos)

		up.ID, err = bkt.NextSequence()
		if err != nil {
			return err
		}

		src, err := store.marshler.Marshal(up)
		if err != nil {
			return err
		}

		return bkt.Put([]byte(strconv.FormatUint(up.ID, 10)), src)
	})
}

func (store *UsersPhotosStore) Update(up *model.UserPhoto) (err error) {
	if up == nil || up.UserID == 0 || up.PhotoID == 0 {
		return nil
	}

	userPhoto := store.Get(up)
	if userPhoto == nil {
		return store.Add(up)
	}

	if up.ID == 0 {
		up.ID = userPhoto.ID
	}

	if up.CreatedAt == 0 {
		up.CreatedAt = userPhoto.CreatedAt
	}

	if up.UpdatedAt <= userPhoto.UpdatedAt {
		up.UpdatedAt = time.Now().UTC().Unix()
	}

	src, err := store.marshler.Marshal(up)
	if err != nil {
		return err
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketUsersPhotos).Put([]byte(strconv.FormatUint(up.ID, 10)), src)
	})
}

func (store *UsersPhotosStore) Get(up *model.UserPhoto) *model.UserPhoto {
	if up == nil || up.UserID == 0 || up.PhotoID == 0 {
		return nil
	}

	userPhoto := new(model.UserPhoto)
	if err := store.conn.View(func(tx *bolt.Tx) (err error) {
		if err = tx.Bucket(common.BucketUsersPhotos).ForEach(func(key, val []byte) (err error) {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if v.GetUint64("user_id") != up.UserID || v.GetUint64("photo_id") != up.PhotoID {
				return nil
			}

			if err = store.marshler.Unmarshal(val, userPhoto); err != nil {
				return err
			}

			return ErrForEachStop
		}); err != nil && !xerrors.Is(err, ErrForEachStop) {
			return err
		}

		return nil
	}); err != nil || userPhoto.PhotoID == 0 || userPhoto.UserID == 0 {
		return nil
	}

	return userPhoto
}

func (store *UsersPhotosStore) GetList(uid uint64, offset, limit int, query string) (model.Photos, int) {
	if limit <= 0 {
		limit = 0
	}

	count := 0
	photos := make(model.Photos, 0, limit)
	_ = store.conn.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(common.BucketPhotos)
		return tx.Bucket(common.BucketUsersPhotos).ForEach(func(key, val []byte) error {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if v.GetUint64("user_id") != uid {
				return nil
			}

			if query != "" && !strings.ContainsAny(string(v.GetStringBytes("query")), query) {
				return nil
			}

			count++

			if (offset != 0 && count <= offset) || (limit > 0 && count > offset+limit) {
				return nil
			}

			p := new(model.Photo)
			if err = store.marshler.Unmarshal(
				bkt.Get([]byte(strconv.FormatUint(v.GetUint64("photo_id"), 10))), p,
			); err != nil {
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

func (store *UsersPhotosStore) Remove(up *model.UserPhoto) (err error) {
	userPhoto := store.Get(up)
	if userPhoto == nil {
		return ErrUserPhotoNotExist
	}

	return store.conn.Update(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket(common.BucketUsersPhotos)
		if err = bkt.ForEach(func(key, val []byte) (err error) {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if v.GetUint64("user_id") != up.UserID || v.GetUint64("photo_id") != up.PhotoID {
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
