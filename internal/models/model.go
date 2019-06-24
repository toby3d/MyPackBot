//go:generate ffjson $GOFILE
package models

type Model struct {
	ID      string `json:"id"`
	SavedAt int64  `json:"saved_at"`
}
