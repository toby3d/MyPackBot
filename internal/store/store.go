package store

import (
	"errors"
	"sort"
	"time"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/photos"
	"gitlab.com/toby3d/mypackbot/internal/model/stickers"
	usersphotos "gitlab.com/toby3d/mypackbot/internal/model/users/photos"
	usersstickers "gitlab.com/toby3d/mypackbot/internal/model/users/stickers"
)

type Store struct {
	Stickers      stickers.Manager
	UsersStickers usersstickers.Manager
	Photos        photos.Manager
	UsersPhotos   usersphotos.Manager
}

// ErrForEachStop used in ForEach loops in database for forse stop iterations
var ErrForEachStop = errors.New("for each stop stop")

func (store *Store) GetList(uid uint64, offset, limit int, query string) (results []interface{}, count int) {
	photos, photosCount := store.UsersPhotos.GetList(uid, 0, -1, query)
	stickers, stickersCount := store.UsersStickers.GetList(uid, 0, -1, query)
	count = stickersCount + photosCount
	results = make([]interface{}, 0)

	for i := range stickers {
		results = append(results, stickers[i])
	}

	for i := range photos {
		results = append(results, photos[i])
	}

	sort.Slice(results, func(i, j int) bool {
		var a, b time.Time

		switch result := results[i].(type) {
		case *model.UserSticker:
			a = time.Unix(result.CreatedAt, 0)
		case *model.UserPhoto:
			a = time.Unix(result.CreatedAt, 0)
		}

		switch result := results[j].(type) {
		case *model.UserSticker:
			b = time.Unix(result.CreatedAt, 0)
		case *model.UserPhoto:
			b = time.Unix(result.CreatedAt, 0)
		}

		return a.Before(b)
	})

	if tail := len(results[offset:]); tail < limit {
		results = results[offset : offset+tail]
	} else {
		results = results[offset : offset+limit-1]
	}

	return results, count
}
