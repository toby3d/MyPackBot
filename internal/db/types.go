package db

import (
	bolt "github.com/etcd-io/bbolt"
	"golang.org/x/text/language"
)

type (
	DataSource interface {
		CreateUser(*User) error
		CreateSet(*Set) error
		CreateSticker(*Sticker) error
		GetUser(*User) (*User, error)
		GetSet(*Set) (*Set, error)
		GetSticker(*Sticker) (*Sticker, error)
		UpdateUser(*User) error
		UpdateSet(*Set) error
		UpdateSticker(*Sticker) error
		DeleteUser(*User) error
		DeleteSet(*Set) error
		DeleteSticker(*Sticker) error

		Close() error
	}

	DB struct {
		db   *bolt.DB
		path string
	}

	User struct {
		Hits     int
		ID       int
		Language language.Tag
		Sets     []*Set
		Sort     string // sort sets in feed
		State    string
	}

	Set struct {
		Hits       int
		ID         string
		IsFavorite bool
		Sort       string // sort stickers in set
		Stickers   []*Sticker
		User       *User
	}

	Sticker struct {
		Emoji      string
		Hits       int
		ID         string
		IsFavorite bool
		Set        *Set
		User       *User
	}
)

const (
	SortAsc  = "asc"
	SortDesc = "desc"
	SortHits = "hits"

	StateNone = "none"
)

var (
	ErrAlreadyExist   = bolt.ErrBucketExists
	ErrDatabaseClosed = bolt.ErrDatabaseNotOpen
	ErrNotFound       = bolt.ErrBucketNotFound

	bucketUsers = []byte("users")

	keyEmoji      = []byte("emoji")
	keyHits       = []byte("hits")
	keyIsFavorite = []byte("is_favorite")
	keyLanguage   = []byte("language")
	keySort       = []byte("sort")
	keyState      = []byte("state")

	valSortAsc  = []byte(SortAsc)
	valSortDesc = []byte(SortDesc)
	valSortHits = []byte(SortHits)
)
