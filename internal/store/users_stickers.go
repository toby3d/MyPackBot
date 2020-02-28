package store

import (
	"strings"
	"time"

	"github.com/timshannon/bolthold"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/stickers"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/xerrors"
)

type UsersStickersStore struct {
	conn     *bolthold.Store
	stickers stickers.Reader
	users    users.Reader
}

func NewUsersStickersStore(conn *bolthold.Store, us users.Reader, ss stickers.Reader) *UsersStickersStore {
	return &UsersStickersStore{
		conn:     conn,
		stickers: ss,
		users:    us,
	}
}

func (store *UsersStickersStore) Add(us *model.UserSticker) error {
	if store.users.Get(us.UserID) == nil || store.stickers.Get(us.StickerID) == nil {
		return bolthold.ErrNotFound
	}

	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.InsertIntoBucket(tx.Bucket(common.BucketUsersStickers), bolthold.NextSequence(), us)
	})
}

func (store *UsersStickersStore) AddSet(uid int, setName string) error {
	for _, sticker := range store.stickers.GetSet(setName) {
		us := model.UserSticker{UserID: uid, StickerID: sticker.ID}

		if store.Get(&us) != nil {
			continue
		}

		now := time.Now().UTC().Unix()
		us.Query = sticker.Emoji
		us.CreatedAt = now
		us.UpdatedAt = now

		if err := store.Add(&us); err != nil {
			return err
		}
	}

	return nil
}

func (store *UsersStickersStore) Update(us *model.UserSticker) (err error) {
	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.UpdateMatchingInBucket(
			tx.Bucket(common.BucketUsersStickers), &model.UserSticker{},
			bolthold.Where("UserID").Eq(us.UserID).And("StickerID").Eq(us.StickerID),
			func(record interface{}) error {
				result, ok := record.(*model.UserSticker) // record will always be a pointer
				if !ok {
					return xerrors.New("invalid type")
				}

				result.Query, result.UpdatedAt = us.Query, us.UpdatedAt

				return nil
			})
	})
}

func (store *UsersStickersStore) Get(us *model.UserSticker) *model.Sticker {
	result := new(model.UserSticker)

	if err := store.conn.Bolt().View(func(tx *bolt.Tx) error {
		return store.conn.FindOneInBucket(
			tx.Bucket(common.BucketUsersStickers), result,
			bolthold.Where("UserID").Eq(us.UserID).And("StickerID").Eq(us.StickerID),
		)
	}); err != nil {
		return nil
	}

	return store.stickers.Get(result.StickerID)
}

func (store *UsersStickersStore) GetList(offset, limit int, filter *model.UserSticker) (list model.Stickers, count int,
	err error) {

	q := bolthold.Where("UserID").Ne(0).And("StickerID").Ne("")
	qCount := bolthold.Where("UserID").Ne(0).And("StickerID").Ne("")

	if offset != 0 {
		q = q.Skip(offset)
	}

	if limit != 0 {
		q = q.Limit(limit)
	}

	if filter != nil {
		if filter.UserID != 0 {
			q = q.And("UserID").Eq(filter.UserID)
			qCount = qCount.And("UserID").Eq(filter.UserID)
		}

		if filter.Query != "" {
			q = q.And("Query").MatchFunc(func(field string) (bool, error) {
				return strings.ContainsAny(field, filter.Query), nil
			})
			qCount = qCount.And("Query").MatchFunc(func(field string) (bool, error) {
				return strings.ContainsAny(field, filter.Query), nil
			})
		}
	}

	results := make(model.UserStickers, 0)

	if err = store.conn.Bolt().View(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket(common.BucketUsersStickers)

		if count, err = store.conn.CountInBucket(bkt, &model.UserSticker{}, qCount); err != nil {
			return err
		}

		return store.conn.FindInBucket(bkt, &results, q)
	}); err != nil {
		return nil, 0, err
	}

	list = make(model.Stickers, 0)

	for i := range results {
		list = append(list, store.stickers.Get(results[i].StickerID))
	}

	return list, count, err
}

func (store *UsersStickersStore) Remove(us *model.UserSticker) error {
	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.DeleteMatchingFromBucket(
			tx.Bucket(common.BucketUsersStickers), &model.UserSticker{},
			bolthold.Where("UserID").Eq(us.UserID).And("StickerID").Eq(us.StickerID),
		)
	})
}

func (store *UsersStickersStore) RemoveSet(uid int, setName string) error {
	ids := make([]string, 0)

	for _, sticker := range store.stickers.GetSet(setName) {
		ids = append(ids, sticker.ID)
	}

	return store.conn.Bolt().Update(func(tx *bolt.Tx) error {
		return store.conn.DeleteMatchingFromBucket(
			tx.Bucket(common.BucketUsersStickers), &model.UserSticker{},
			bolthold.Where("UserID").Eq(uid).And("StickerID").ContainsAny(ids),
		)
	})
}
