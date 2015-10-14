package term

import (
	"os"

	"github.com/pkg/term"
)

type Term struct {
	*term.Term
}

func RawMode(t *Term) error {
	return term.RawMode(t.Term)
}

func New(f *os.File, options ...func(*Term) error) (*Term, error) {
	t, err := term.Open(f.Name())
	if err != nil {
		return nil, err
	}
	tw := &Term{Term: t}
	if err := tw.SetOption(options...); err != nil {
		return nil, err
	}
	return tw, nil
}

func (t *Term) SetOption(options ...func(*Term) error) error {
	for _, opt := range options {
		if err := opt(t); err != nil {
			return err
		}
	}
	return nil
}
