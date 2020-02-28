package store

import (
	"strings"

	"github.com/timshannon/bolthold"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/photos"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/xerrors"
)

type UsersPhotosStore struct {
	conn   *bolthold.Store
	photos photos.Reader
	users  users.Reader
}

func NewUsersPhotosStore(conn *bolthold.Store, us users.Reader, ps photos.Reader) *UsersPhotosStore {
	return &UsersPhotosStore{
		conn:   conn,
		photos: ps,
		users:  us,
	}
}

func (store *UsersPhotosStore) Add(up *model.UserPhoto) (err error) {
	if store.users.Get(up.UserID) == nil || store.photos.Get(up.PhotoID) == nil {
		return bolthold.ErrNotFound
	}

	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.InsertIntoBucket(tx.Bucket(common.BucketUsersPhotos), bolthold.NextSequence(), up)
	})
}

func (store *UsersPhotosStore) Update(up *model.UserPhoto) (err error) {
	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.UpdateMatchingInBucket(
			tx.Bucket(common.BucketUsersPhotos), &model.UserPhoto{},
			bolthold.Where("UserID").Eq(up.UserID).And("PhotoID").Eq(up.PhotoID),
			func(record interface{}) error {
				result, ok := record.(*model.UserPhoto) // record will always be a pointer
				if !ok {
					return xerrors.New("invalid type")
				}

				result.Query, result.UpdatedAt = up.Query, up.UpdatedAt

				return nil
			})
	})

}

func (store *UsersPhotosStore) Get(up *model.UserPhoto) *model.Photo {
	result := new(model.UserPhoto)

	if err := store.conn.Bolt().View(func(tx *bolt.Tx) error {
		return store.conn.FindOneInBucket(
			tx.Bucket(common.BucketUsersPhotos), result,
			bolthold.Where("UserID").Eq(up.UserID).And("PhotoID").Eq(up.PhotoID),
		)
	}); err != nil {
		return nil
	}

	return store.photos.Get(result.PhotoID)
}

func (store *UsersPhotosStore) GetList(offset, limit int, filter *model.UserPhoto) (list model.Photos, count int,
	err error) {

	q := bolthold.Where("UserID").Ne(0).And("PhotoID").Ne("")
	qCount := bolthold.Where("UserID").Ne(0).And("PhotoID").Ne("")

	if offset != 0 {
		q = q.Skip(offset)
	}

	if limit != 0 {
		q = q.Limit(limit)
	}

	if filter != nil {
		if filter.UserID != 0 {
			q = q.And("UserID").Eq(filter.UserID)
			qCount.And("UserID").Eq(filter.UserID)
		}

		if filter.Query != "" {
			q = q.And("Query").MatchFunc(func(field string) (bool, error) {
				return strings.ContainsAny(field, filter.Query), nil
			})
			qCount.And("Query").MatchFunc(func(field string) (bool, error) {
				return strings.ContainsAny(field, filter.Query), nil
			})
		}
	}

	results := make(model.UserPhotos, 0)

	if err = store.conn.Bolt().View(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket(common.BucketUsersPhotos)

		if count, err = store.conn.CountInBucket(bkt, &model.UserPhoto{}, qCount); err != nil {
			return err
		}

		return store.conn.FindInBucket(bkt, &results, q)
	}); err != nil {
		return nil, 0, err
	}

	list = make(model.Photos, 0)

	for i := range results {
		list = append(list, store.photos.Get(results[i].PhotoID))
	}

	return list, count, err
}

func (store *UsersPhotosStore) Remove(up *model.UserPhoto) (err error) {
	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.DeleteMatchingFromBucket(
			tx.Bucket(common.BucketUsersPhotos), &model.UserPhoto{},
			bolthold.Where("UserID").Eq(up.UserID).And("PhotoID").Eq(up.PhotoID),
		)
	})
}
