package model

import (
	"fmt"

	"golang.org/x/xerrors"
)

type Error struct {
	Message string
	Frame   xerrors.Frame
}

func (err Error) Error() string {
	return fmt.Sprint(err)
}

func (err Error) Format(f fmt.State, c rune) {
	xerrors.FormatError(err, f, c)
}

func (err Error) FormatError(p xerrors.Printer) error {
	p.Print(err.Message)

	if p.Detail() {
		err.Frame.Format(p)
	}

	return nil
}
