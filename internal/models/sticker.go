//go:generate ffjson $GOFILE
package models

type Sticker struct {
	ID      string `json:"id"`
	Emoji   string `json:"emoji"`
	SetName string `json:"set_name"`
}
