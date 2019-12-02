//nolint: gochecknoglobals
package common

import (
	"github.com/Masterminds/semver"
	tg "gitlab.com/toby3d/telegram"
)

const (
	CommandEdit string = "edit"

	// NOTE(toby3d): DEPRECATED
	CommandAddPack    string = "addpack"
	CommandAddSticker string = "add" + tg.TypeSticker
	CommandCancel     string = "cancel"
	CommandDelPack    string = "addpack"
	CommandDelSticker string = "add" + tg.TypeSticker
	CommandReset      string = "reset"
)

const (
	DataSeparator string = "@"
	DataAdd       string = "add"
	DataDel       string = "del"
	DataSet       string = "set"

	DataAddSet string = DataAdd + DataSeparator + DataSet
	DataDelSet string = DataDel + DataSeparator + DataSet
)

const SetNameUploaded string = "uploaded_by_mypackbot"

var Version = semver.MustParse("2.0.0")

var (
	BucketPhotos        = []byte("photos")
	BucketStickers      = []byte("stickers")
	BucketUsers         = []byte("users")
	BucketUsersPhotos   = []byte("users_photos")
	BucketUsersStickers = []byte("users_stickers")
	Buckets             = [...][]byte{
		BucketPhotos,
		BucketStickers,
		BucketUsers,
		BucketUsersPhotos,
		BucketUsersStickers,
	}
)
