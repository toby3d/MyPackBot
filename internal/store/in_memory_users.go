package store

import (
	"sort"
	"sync"
	"time"

	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
)

type InMemoryUsersStore struct {
	mutex sync.RWMutex
	users model.Users
}

func NewInMemoryUsersStore() *InMemoryUsersStore {
	return &InMemoryUsersStore{
		mutex: sync.RWMutex{},
		users: make([]*model.User, 0),
	}
}

func (store *InMemoryUsersStore) Create(u *model.User) error {
	if store.Get(u.ID) != nil {
		return common.ErrUserExist
	}

	if u.CreatedAt == 0 {
		u.CreatedAt = time.Now().UTC().Unix()
	}

	store.mutex.Lock()
	store.users = append(store.users, u)

	sort.Slice(store.users, func(i, j int) bool { return store.users[i].ID < store.users[j].ID })
	store.mutex.Unlock()

	return nil
}

func (store *InMemoryUsersStore) Get(uid int) *model.User {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	return store.users.GetByID(uid)
}

func (store *InMemoryUsersStore) Update(u *model.User) error {
	if store.Get(u.ID) == nil {
		return store.Create(u)
	}

	if u.UpdatedAt <= 0 {
		u.UpdatedAt = time.Now().UTC().Unix()
	}

	store.mutex.Lock()
	for i := range store.users {
		if store.users[i].ID != u.ID {
			continue
		}

		store.users[i] = u
	}
	store.mutex.Unlock()

	return nil
}

func (store *InMemoryUsersStore) Remove(uid int) error {
	if store.Get(uid) == nil {
		return common.ErrUserNotExist
	}

	store.mutex.Lock()
	for i := range store.users {
		if store.users[i].ID != uid {
			continue
		}

		store.users = store.users[:i+copy(store.users[i:], store.users[i+1:])]

		break
	}

	sort.Slice(store.users, func(i, j int) bool {
		return store.users[i].ID < store.users[j].ID
	})
	store.mutex.Unlock()

	return nil
}

func (store *InMemoryUsersStore) GetOrCreate(u *model.User) (*model.User, error) {
	if user := store.Get(u.ID); user != nil {
		return user, nil
	}

	if err := store.Create(u); err != nil {
		return nil, err
	}

	return store.Get(u.ID), nil
}
