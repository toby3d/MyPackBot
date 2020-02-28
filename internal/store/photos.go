package store

import (
	"github.com/timshannon/bolthold"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	bolt "go.etcd.io/bbolt"
)

type PhotosStore struct {
	conn *bolthold.Store
}

func NewPhotosStore(conn *bolthold.Store) *PhotosStore {
	return &PhotosStore{conn: conn}
}

func (store *PhotosStore) Create(p *model.Photo) error {
	if store.Get(p.ID) != nil {
		return bolthold.ErrKeyExists
	}

	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.InsertIntoBucket(tx.Bucket(common.BucketPhotos), p.ID, p)
	})
}

func (store *PhotosStore) Get(id string) *model.Photo {
	result := new(model.Photo)

	if err := store.conn.Bolt().View(func(tx *bolt.Tx) error {
		return store.conn.GetFromBucket(tx.Bucket(common.BucketPhotos), id, result)
	}); err != nil {
		return nil
	}

	return result
}

func (store *PhotosStore) GetList(offset, limit int, filter *model.Photo) (list model.Photos, count int, err error) {
	q := bolthold.Where("ID").Ne("")

	if offset > 0 {
		q = q.Skip(offset)
	}

	if limit > 0 {
		q = q.Limit(limit)
	}

	// TODO(toby3d): implement filter here
	list = make(model.Photos, limit)

	if err := store.conn.Bolt().View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(common.BucketPhotos)

		if count, err = store.conn.CountInBucket(bkt, &model.Photo{}, q); err != nil {
			return err
		}

		return store.conn.FindInBucket(bkt, &list, q)
	}); err != nil {
		return list, count, err
	}

	return list, count, err
}

func (store *PhotosStore) Update(p *model.Photo) error {
	if store.Get(p.ID) == nil {
		return store.Create(p)
	}

	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.UpdateBucket(tx.Bucket(common.BucketPhotos), p.ID, p)
	})
}

func (store *PhotosStore) Remove(id string) error {
	if store.Get(id) == nil {
		return bolthold.ErrNotFound
	}

	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.DeleteFromBucket(tx.Bucket(common.BucketPhotos), id, &model.Photo{})
	})
}

func (store *PhotosStore) GetOrCreate(p *model.Photo) (*model.Photo, error) {
	if photo := store.Get(p.ID); photo != nil {
		return photo, nil
	}

	if err := store.Create(p); err != nil {
		return nil, err
	}

	return store.GetOrCreate(p)
}
