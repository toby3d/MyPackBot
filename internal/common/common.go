//nolint: gochecknoglobals
package common

import (
	"github.com/Masterminds/semver"
)

const (
	CommandPing       string = "ping"
	CommandAddSticker string = "addsticker"
	CommandAddPack    string = "addpack"
	CommandDelSticker string = "delsticker"
	CommandDelPack    string = "delpack"
	CommandReset      string = "reset"
	CommandCancel     string = "cancel"
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
