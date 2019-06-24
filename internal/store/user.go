package store

import (
	"strconv"
	"time"

	bolt "github.com/etcd-io/bbolt"
	json "github.com/pquerna/ffjson/ffjson"
	"gitlab.com/toby3d/mypackbot/internal/models"
)

type UserStore struct {
	db *bolt.DB
}

var (
	bktUsers         = []byte("users")
	bktUsersStickers = []byte("users_stickers")
)

func NewUserStore(db *bolt.DB) *UserStore {
	return &UserStore{db: db}
}

func (us *UserStore) GetByID(uid int) (*models.User, error) {
	var u models.User
	err := us.db.View(func(tx *bolt.Tx) error {
		src := tx.Bucket(bktUsers).Get([]byte(strconv.Itoa(uid)))
		if src == nil {
			return nil
		}
		return json.UnmarshalFast(src, &u)
	})
	return &u, err
}

func (us *UserStore) Create(u *models.User) error {
	if u.ID == 0 {
		return nil
	}

	u.StartedAt = time.Now().UTC().UnixNano()
	src, err := json.MarshalFast(u)
	if err != nil {
		return err
	}

	return us.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bktUsers).Put([]byte(strconv.Itoa(u.ID)), src)
	})
}

func (us *UserStore) Update(u *models.User) error {
	if u.ID == 0 {
		return nil
	}

	src, err := json.MarshalFast(u)
	if err != nil {
		return err
	}

	return us.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bktUsers).Put([]byte(strconv.Itoa(u.ID)), src)
	})
}

func (us *UserStore) AddSticker(uid int, sid string) error {
	src, err := json.MarshalFast(&models.UsersStickers{
		UserID:    uid,
		StickerID: sid,
	})
	if err != nil {
		return err
	}

	return us.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktUsersStickers)
		id, err := bkt.NextSequence()
		if err != nil {
			return err
		}

		return bkt.Put(strconv.AppendUint(nil, id, 10), src)
	})
}

func (us *UserStore) DeleteSticker(uid int, sid string) error {
	return us.db.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket(bktUsersStickers).Cursor()
		for _, v := c.First(); v != nil; _, v = c.Next() {
			var us models.UsersStickers
			if err := json.UnmarshalFast(v, &us); err != nil {
				return err
			}

			if us.UserID != uid || us.StickerID != sid {
				continue
			}
			return c.Delete()
		}
		return nil
	})
}
