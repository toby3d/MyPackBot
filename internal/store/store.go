package store

import (
	"sort"
	"strconv"
	"strings"

	bolt "github.com/etcd-io/bbolt"
	json "github.com/pquerna/ffjson/ffjson"
	"gitlab.com/toby3d/mypackbot/internal/models"
)

type Store struct {
	db            *bolt.DB
	stickers      []models.Sticker
	users         []models.User
	usersStickers []models.UsersStickers
	sets          []models.Set
}

var (
	bktSets          = []byte("sets")
	bktStickers      = []byte("stickers")
	bktUsers         = []byte("users")
	bktUsersStickers = []byte("users_stickers")
)

func New(dataBase *bolt.DB) (*Store, error) {
	var store Store
	store.db = dataBase
	err := store.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) (err error) {
			switch string(name) {
			case string(bktUsers):
				return b.ForEach(func(k, v []byte) error {
					var u models.User
					if err := json.UnmarshalFast(v, &u); err != nil {
						return err
					}
					store.users = append(store.users, u)
					return nil
				})
			case string(bktStickers):
				return b.ForEach(func(k, v []byte) error {
					var s models.Sticker
					if err := json.UnmarshalFast(v, &s); err != nil {
						return err
					}
					store.stickers = append(store.stickers, s)
					return nil
				})
			case string(bktUsersStickers):
				return b.ForEach(func(k, v []byte) error {
					var us models.UsersStickers
					if err := json.UnmarshalFast(v, &us); err != nil {
						return err
					}
					store.usersStickers = append(store.usersStickers, us)
					return nil
				})
			case string(bktSets):
				return b.ForEach(func(k, v []byte) error {
					var s models.Set
					if err := json.UnmarshalFast(v, &s); err != nil {
						return err
					}
					store.sets = append(store.sets, s)
					return nil
				})
			default:
				return nil
			}
		})
	})
	return &store, err
}

func (s *Store) GetUser(uid int) *models.User {
	for _, u := range s.users {
		u := u
		if u.ID != uid {
			continue
		}
		return &u
	}
	return nil
}

func (s *Store) GetOrCreateUser(src *models.User) (*models.User, error) {
	if user := s.GetUser(src.ID); user != nil {
		return user, nil
	}
	if err := s.CreateUser(src); err != nil {
		return nil, err
	}
	return s.GetUser(src.ID), nil
}

func (s *Store) GetSet(setName string) *models.Set {
	for _, set := range s.sets {
		set := set
		if set.Name != setName {
			continue
		}
		return &set
	}
	return nil
}

func (s *Store) GetOrCreateSet(src *models.Set) (*models.Set, error) {
	if set := s.GetSet(src.Name); set != nil {
		return set, nil
	}
	if err := s.CreateSet(src); err != nil {
		return nil, err
	}
	return s.GetSet(src.Name), nil
}

func (s *Store) GetSticker(sid string) *models.Sticker {
	for _, s := range s.stickers {
		s := s
		if s.ID != sid {
			continue
		}
		return &s
	}
	return nil
}

func (s *Store) GetOrCreateSticker(src *models.Sticker) (*models.Sticker, error) {
	if sticker := s.GetSticker(src.ID); sticker != nil {
		return sticker, nil
	}
	if err := s.CreateSticker(src); err != nil {
		return nil, err
	}
	return s.GetSticker(src.ID), nil
}

func (s *Store) GetUsers(offset, limit int) (users []models.User, count int) {
	for _, u := range s.users {
		u := u
		count++
		if offset != -1 && offset > 0 {
			offset--
			continue
		}
		if limit != -1 && limit > 0 {
			limit--
		}
		if limit == 0 {
			continue
		}

		users = append(users, u)
	}
	return
}

func (s *Store) GetStickers(offset, limit int) (stickers []models.Sticker, count int) {
	for _, s := range s.stickers {
		s := s
		count++
		if offset != -1 && offset > 0 {
			offset--
			continue
		}
		if limit != -1 && limit > 0 {
			limit--
		}
		if limit == 0 {
			continue
		}

		stickers = append(stickers, s)
	}
	return
}

func (s *Store) GetStickersBySet(setName string, offset, limit int) (stickers []models.Sticker, count int) {
	for _, s := range s.stickers {
		s := s
		if s.SetName != setName {
			continue
		}
		count++

		if offset != -1 && offset > 0 {
			offset--
			continue
		}

		if limit != -1 && limit > 0 {
			limit--
		}
		if limit == 0 {
			continue
		}

		stickers = append(stickers, s)
	}
	return
}

func (s *Store) GetUserStickers(uid, offset, limit int) (stickers []models.Sticker, count int) {
	for _, us := range s.usersStickers {
		us := us
		if us.UserID != uid {
			continue
		}
		count++

		if offset != -1 && offset > 0 {
			offset--
			continue
		}

		if limit != -1 && limit > 0 {
			limit--
		}
		if limit == 0 {
			continue
		}

		sticker := s.GetSticker(us.StickerID)
		stickers = append(stickers, *sticker)
	}
	return
}

func (s *Store) GetUserStickersByQuery(query string, uid, offset, limit int) (stickers []models.Sticker, count int) {
	userStickers, _ := s.GetUserStickers(uid, -1, -1)
	for _, s := range userStickers {
		s := s
		if !strings.ContainsAny(s.Emoji, query) {
			continue
		}
		count++

		if offset != -1 && offset > 0 {
			offset--
			continue
		}

		if limit != -1 && limit > 0 {
			limit--
		}
		if limit == 0 {
			continue
		}

		stickers = append(stickers, s)
	}
	return
}

func (s *Store) GetUserStickersBySet(setName string, uid, offset, limit int) (stickers []models.Sticker, count int) {
	userStickers, _ := s.GetUserStickers(uid, -1, -1)
	for _, s := range userStickers {
		s := s
		if s.SetName != setName {
			continue
		}
		count++

		if offset != -1 && offset > 0 {
			offset--
			continue
		}

		if limit != -1 && limit > 0 {
			limit--
		}
		if limit == 0 {
			continue
		}

		stickers = append(stickers, s)
	}
	return
}

func (s *Store) CreateSet(set *models.Set) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		tx.Rollback()
		return err
	}

	if s.GetSet(set.Name) != nil {
		return tx.Rollback()
	}

	src, err := json.MarshalFast(set)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Bucket(bktSets).Put([]byte(set.Name), src); err != nil {
		tx.Rollback()
		return err
	}

	s.sets = append(s.sets, *set)
	return tx.Commit()
}

func (s *Store) CreateUser(user *models.User) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		tx.Rollback()
		return err
	}

	if s.GetUser(user.ID) != nil {
		return tx.Rollback()
	}

	src, err := json.MarshalFast(user)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Bucket(bktUsers).Put([]byte(strconv.Itoa(user.ID)), src); err != nil {
		tx.Rollback()
		return err
	}

	s.users = append(s.users, *user)
	return tx.Commit()
}

func (s *Store) CreateSticker(sticker *models.Sticker) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		tx.Rollback()
		return err
	}

	if s.GetSticker(sticker.ID) != nil {
		return tx.Rollback()
	}

	src, err := json.MarshalFast(sticker)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Bucket(bktStickers).Put([]byte(sticker.ID), src); err != nil {
		tx.Rollback()
		return err
	}

	s.stickers = append(s.stickers, *sticker)
	sort.Slice(s.stickers, func(i, j int) bool {
		return s.stickers[i].SetName != s.stickers[j].SetName
	})
	return tx.Commit()
}

func (s *Store) UpdateUser(user *models.User) error {
	return nil
}

func (s *Store) UpdateSticker(sticker *models.Sticker) error {
	return nil
}

func (s *Store) UpdateSet(user *models.User) error {
	return nil
}

func (s *Store) DeleteUser(uid int) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		tx.Rollback()
		return err
	}

	if s.GetUser(uid) == nil {
		return tx.Rollback()
	}

	if err = tx.Bucket(bktUsers).Delete([]byte(strconv.Itoa(uid))); err != nil {
		tx.Rollback()
		return err
	}

	userStickers := tx.Bucket(bktUsersStickers)
	c := userStickers.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var us models.UsersStickers
		if err = json.UnmarshalFast(v, &us); err != nil {
			tx.Rollback()
			return err
		}
		if us.UserID != uid {
			continue
		}
		if err = userStickers.Delete(k); err != nil {
			tx.Rollback()
			return err
		}
	}

	for i := range s.users {
		if s.users[i].ID != uid {
			continue
		}
		s.users = append(s.users[:i], s.users[i+1:]...)
	}

	for i := range s.usersStickers {
		if s.usersStickers[i].UserID != uid {
			continue
		}
		s.usersStickers = append(s.usersStickers[:i], s.usersStickers[i+1:]...)
	}

	return tx.Commit()
}

func (s *Store) DeleteSticker(sid string) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		tx.Rollback()
		return err
	}

	if s.GetSticker(sid) == nil {
		tx.Rollback()
		return nil
	}

	if err = tx.Bucket(bktStickers).Delete([]byte(sid)); err != nil {
		tx.Rollback()
		return err
	}

	c := tx.Bucket(bktUsersStickers).Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var us models.UsersStickers
		if err = json.UnmarshalFast(v, &us); err != nil {
			tx.Rollback()
			return err
		}
		if us.StickerID != sid {
			continue
		}
		if err = c.Delete(); err != nil {
			tx.Rollback()
			return err
		}
	}

	for i := range s.stickers {
		if s.stickers[i].ID != sid {
			continue
		}
		s.stickers = append(s.stickers[:i], s.stickers[i+1:]...)
	}

	for i := range s.usersStickers {
		if s.usersStickers[i].StickerID != sid {
			continue
		}
		s.usersStickers = append(s.usersStickers[:i], s.usersStickers[i+1:]...)
	}

	return tx.Commit()
}

func (s *Store) DeleteSet(setName string) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		tx.Rollback()
		return err
	}

	if s.GetSet(setName) == nil {
		tx.Rollback()
		return nil
	}

	if err = tx.Bucket(bktSets).Delete([]byte(setName)); err != nil {
		tx.Rollback()
		return err
	}

	stickers, _ := s.GetStickersBySet(setName, -1, -1)
	for _, sticker := range stickers {
		sticker := sticker
		if err = s.DeleteSticker(sticker.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	for i := range s.sets {
		if s.sets[i].Name != setName {
			continue
		}
		s.sets = append(s.sets[:i], s.sets[i+1:]...)
	}

	return tx.Commit()
}

func (s *Store) AddSticker(uid int, sid string) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, us := range s.usersStickers {
		if us.UserID != uid || us.StickerID != sid {
			continue
		}
		tx.Rollback()
		return nil
	}

	src, err := json.MarshalFast(&models.UsersStickers{
		UserID:    uid,
		StickerID: sid,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	usersStickers := tx.Bucket(bktUsersStickers)
	id, err := usersStickers.NextSequence()
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = usersStickers.Put(strconv.AppendUint(nil, id, 10), src); err != nil {
		tx.Rollback()
		return err
	}

	s.usersStickers = append(s.usersStickers, models.UsersStickers{
		UserID:    uid,
		StickerID: sid,
	})
	return tx.Commit()
}

func (s *Store) RemoveSticker(uid int, sid string) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		tx.Rollback()
		return err
	}

	c := tx.Bucket(bktUsersStickers).Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var us models.UsersStickers
		if err = json.UnmarshalFast(v, &us); err != nil {
			tx.Rollback()
			return err
		}
		if us.UserID != uid || us.StickerID != sid {
			continue
		}
		if err = c.Delete(); err != nil {
			tx.Rollback()
			return err
		}
	}

	for i := range s.usersStickers {
		if s.usersStickers[i].UserID == uid && s.usersStickers[i].StickerID == sid {
			s.usersStickers = append(s.usersStickers[:i], s.usersStickers[i+1:]...)
		}
	}

	return tx.Commit()
}
