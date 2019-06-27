//go:generate ffjson $GOFILE
package models

type Sticker struct {
	Model
	Emoji   string `json:"emoji"`
	SetName string `json:"set_name"`
}

func (s *Sticker) InSet() bool { return s.SetName != "" }
