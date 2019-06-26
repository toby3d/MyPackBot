//go:generate ffjson $GOFILE
package models

type Set struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}
