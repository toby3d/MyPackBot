//nolint: gochecknoglobals
package common

import (
	"github.com/Masterminds/semver"
	"gitlab.com/toby3d/mypackbot/internal/model"
)

const (
	CommandPing string = "ping"
)

const (
	DataAdd           string = "add"
	DataAddSet        string = DataSet + DataSeparator + DataAdd
	DataAddSticker    string = DataSticker + DataSeparator + DataAdd
	DataLanguage      string = "language"
	DataRemove        string = "remove"
	DataRemoveSet     string = DataSet + DataSeparator + DataRemove
	DataRemoveSticker string = DataSticker + DataSeparator + DataRemove
	DataSeparator     string = "@"
	DataSet           string = "set"
	DataSticker       string = "sticker"
)

const (
	SetNameUploaded string = "uploaded_by_mypackbot"
)

var Version = semver.MustParse("2.0.0")

var (
	BucketStickers      = []byte("stickers")
	BucketUsers         = []byte("users")
	BucketUsersStickers = []byte("users_stickers")
	Buckets             = [...][]byte{BucketStickers, BucketUsers, BucketUsersStickers}
)

var (
	ErrStickerExist        = model.Error{Message: "Sticker already exist"}
	ErrStickerNotExist     = model.Error{Message: "Sticker not exist"}
	ErrUserExist           = model.Error{Message: "User already exist"}
	ErrUserNotExist        = model.Error{Message: "User not exist"}
	ErrUserStickerExist    = model.Error{Message: "Sticker already imported"}
	ErrUserStickerNotExist = model.Error{Message: "Sticker already removed"}
)
