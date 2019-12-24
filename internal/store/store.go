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
	tg "gitlab.com/toby3d/telegram"
)

type (
	Store struct {
		Stickers      stickers.Manager
		UsersStickers usersstickers.Manager
		Photos        photos.Manager
		UsersPhotos   usersphotos.Manager
	}

	Filter struct {
		AllowedTypes []string
		Query        string
		Offset       int
		Limit        int
		IsPersonal   bool
	}
)

// ErrForEachStop used in ForEach loops in database for forse stop iterations
var ErrForEachStop = errors.New("for each stop stop")

func (store *Store) GetList(uid uint64, f *Filter) ([]interface{}, int) {
	results := make([]interface{}, 0)
	count := 0

	for _, t := range f.AllowedTypes {
		switch t {
		case tg.TypePhoto:
			photos, photosCount := store.UsersPhotos.GetList(uid, 0, -1, f.Query)
			for i := range photos {
				results = append(results, photos[i])
			}

			count += photosCount
		case tg.TypeSticker:
			stickers, stickersCount := store.UsersStickers.GetList(uid, 0, -1, f.Query)
			for i := range stickers {
				results = append(results, stickers[i])
			}

			count += stickersCount
		}
	}

	sort.Slice(results, func(i, j int) bool {
		var a, b time.Time

		switch result := results[i].(type) {
		case *model.Sticker:
			a = time.Unix(result.CreatedAt, 0)
		case *model.Photo:
			a = time.Unix(result.CreatedAt, 0)
		}

		switch result := results[j].(type) {
		case *model.Sticker:
			b = time.Unix(result.CreatedAt, 0)
		case *model.Photo:
			b = time.Unix(result.CreatedAt, 0)
		}

		return a.Before(b)
	})

	if len(results) <= f.Offset {
		return make([]interface{}, 0), count
	}

	if tail := len(results[f.Offset:]); tail < f.Limit {
		return results[f.Offset : f.Offset+tail], count
	}

	return results[f.Offset : f.Offset+f.Limit-1], count
}
