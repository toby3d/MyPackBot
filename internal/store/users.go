package store

import (
	"github.com/timshannon/bolthold"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	bolt "go.etcd.io/bbolt"
)

type UsersStore struct {
	conn *bolthold.Store
}

func NewUsersStore(conn *bolthold.Store) *UsersStore {
	return &UsersStore{conn: conn}
}

func (store *UsersStore) Create(u *model.User) error {
	if store.Get(u.ID) != nil {
		return bolthold.ErrKeyExists
	}

	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.InsertIntoBucket(tx.Bucket(common.BucketUsers), u.ID, u)
	})
}

func (store *UsersStore) Get(id int) *model.User {
	result := new(model.User)

	if err := store.conn.Bolt().View(func(tx *bolt.Tx) error {
		return store.conn.GetFromBucket(tx.Bucket(common.BucketUsers), id, result)
	}); err != nil {
		return nil
	}

	return result
}

func (store *UsersStore) GetList(offset, limit int, filter *model.User) (list model.Users, count int, err error) {
	list = make(model.Users, limit)
	q := bolthold.Where("ID").Ne("")

	if offset > 0 {
		q = q.Skip(offset)
	}

	if limit > 0 {
		q = q.Limit(limit)
	}

	// TODO(toby3d): implement filter here

	if err := store.conn.Bolt().View(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket(common.BucketUsers)
		if count, err = store.conn.CountInBucket(bkt, &model.User{}, q); err != nil {
			return err
		}

		return store.conn.FindInBucket(bkt, &list, q)
	}); err != nil {
		return nil, 0, err
	}

	return list, count, err
}

func (store *UsersStore) Update(u *model.User) error {
	if store.Get(u.ID) == nil {
		return store.Create(u)
	}

	return store.conn.Bolt().Update(func(tx *bolt.Tx) (err error) {
		return store.conn.UpdateBucket(tx.Bucket(common.BucketUsers), u.ID, u)
	})
}

func (store *UsersStore) GetOrCreate(u *model.User) (*model.User, error) {
	if user := store.Get(u.ID); user != nil {
		return user, nil
	}

	if err := store.Create(u); err != nil {
		return nil, err
	}

	return u, nil
}
