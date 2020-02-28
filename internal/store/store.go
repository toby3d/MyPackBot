package store

import (
	"sort"

	"gitlab.com/toby3d/mypackbot/internal/model"
	usersphotos "gitlab.com/toby3d/mypackbot/internal/model/users/photos"
	usersstickers "gitlab.com/toby3d/mypackbot/internal/model/users/stickers"
	tg "gitlab.com/toby3d/telegram"
)

type (
	Store struct {
		usersStickers usersstickers.ReadWriter
		usersPhotos   usersphotos.ReadWriter
	}

	Filter struct {
		UserID       int
		AllowedTypes []string
		Query        string
		Offset       int
		Limit        int
	}
)

func (f *Filter) offset() int {
	if len(f.AllowedTypes) == 0 || f.Offset <= len(f.AllowedTypes) {
		return f.Offset
	}

	return f.Offset / len(f.AllowedTypes)
}

func (f *Filter) limit() int {
	if len(f.AllowedTypes) == 0 || f.Limit <= len(f.AllowedTypes) {
		return f.Limit
	}

	return f.Limit / len(f.AllowedTypes)
}

func NewStore(us usersstickers.ReadWriter, up usersphotos.ReadWriter) *Store {
	return &Store{
		usersStickers: us,
		usersPhotos:   up,
	}
}

func (store *Store) GetList(offset, limit int, filter *Filter) (list []model.InlineResult, count int, err error) {
	if filter == nil {
		filter = new(Filter)
		filter.AllowedTypes = []string{tg.TypeSticker, tg.TypePhoto}
	}

	list = make([]model.InlineResult, 0)

	for _, t := range filter.AllowedTypes {
		switch t {
		case tg.TypeSticker:
			l, c, err := store.usersStickers.GetList(
				filter.offset(), filter.limit(), &model.UserSticker{
					UserID: filter.UserID,
					Query:  filter.Query,
				},
			)
			if err != nil {
				return list, count, err
			}

			for i := range l {
				list = append(list, l[i])
			}

			count += c
		case tg.TypePhoto:
			l, c, err := store.usersPhotos.GetList(
				filter.offset(), filter.limit(), &model.UserPhoto{
					UserID: filter.UserID,
					Query:  filter.Query,
				},
			)
			if err != nil {
				return list, count, err
			}

			for i := range l {
				list = append(list, l[i])
			}

			count += c
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].GetUpdatedAt() < list[j].GetUpdatedAt()
	})

	return list, count, err
}
