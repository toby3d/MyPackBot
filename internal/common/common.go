//nolint: gochecknoglobals
package common

import "gitlab.com/toby3d/mypackbot/internal/model"

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

const (
	ContextPrinter model.ContextKey = "printer"
	ContextSticker model.ContextKey = "sticker"
	ContextUser    model.ContextKey = "user"
)

var (
	BucketStickers      = []byte("stickers")
	BucketUsers         = []byte("users")
	BucketUsersStickers = []byte("users_stickers")
)
