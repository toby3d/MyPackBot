package store

import (
	"github.com/timshannon/bolthold"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	bolt "go.etcd.io/bbolt"
)

type StickersStore struct {
	conn *bolthold.Store
}

func NewStickersStore(conn *bolthold.Store) *StickersStore {
	return &StickersStore{conn: conn}
}

func (store *StickersStore) Create(s *model.Sticker) error {
	if store.Get(s.ID) != nil {
		return bolthold.ErrKeyExists
	}

	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.InsertIntoBucket(tx.Bucket(common.BucketStickers), s.ID, s)
	})
}

func (store *StickersStore) Get(id string) *model.Sticker {
	result := new(model.Sticker)

	if err := store.conn.Bolt().View(func(tx *bolt.Tx) error {
		return store.conn.GetFromBucket(tx.Bucket(common.BucketStickers), id, result)
	}); err != nil {
		return nil
	}

	return result
}

func (store *StickersStore) GetSet(name string) model.Stickers {
	list, _, _ := store.GetList(0, 0, &model.Sticker{SetName: name})
	return list
}

func (store *StickersStore) GetList(offset, limit int, filter *model.Sticker) (list model.Stickers, count int,
	err error) {
	q := bolthold.Where("ID").Ne("")

	if offset > 0 {
		q = q.Skip(offset)
	}

	if limit > 0 {
		q = q.Limit(limit)
	}

	if filter != nil {
		if filter.Emoji != "" {
			q = q.And("Emoji").ContainsAny(filter.Emoji)
		}

		if filter.SetName != "" {
			q = q.And("SetName").Eq(filter.SetName)
		}
	}

	list = make(model.Stickers, limit)

	if err := store.conn.Bolt().View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(common.BucketStickers)

		if count, err = store.conn.CountInBucket(bkt, &model.Sticker{}, q); err != nil {
			return err
		}

		return store.conn.FindInBucket(bkt, &list, q)
	}); err != nil {
		return list, count, err
	}

	return list, count, err
}

func (store *StickersStore) Update(s *model.Sticker) error {
	if store.Get(s.ID) == nil {
		return store.Create(s)
	}

	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.UpdateBucket(tx.Bucket(common.BucketStickers), s.ID, s)
	})
}

func (store *StickersStore) Remove(id string) error {
	if store.Get(id) == nil {
		return bolthold.ErrNotFound
	}

	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.DeleteFromBucket(tx.Bucket(common.BucketStickers), id, &model.Sticker{})
	})
}

func (store *StickersStore) GetOrCreate(s *model.Sticker) (*model.Sticker, error) {
	if sticker := store.Get(s.ID); sticker != nil {
		return sticker, nil
	}

	if err := store.Create(s); err != nil {
		return nil, err
	}

	return store.GetOrCreate(s)
}
