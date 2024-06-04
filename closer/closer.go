package closer

import (
	"github.com/rs/zerolog"
	"io"
)

type CloseFunc func() error

func (f CloseFunc) Close() error {
	err := f()
	if err != nil {
		return err
	}

	return nil
}

type C struct {
	Closer io.Closer
	Name   string
}

type Closer struct {
	logger   *zerolog.Logger
	forClose []C
}

func (c *Closer) Close() {
	for _, closer := range c.forClose {
		err := closer.Closer.Close()
		if err != nil {
			c.logger.Err(err).Msgf("failed to close %s: %s", closer.Name, err)
		}
	}

	c.logger.Info().Msg("all dependencies are closed")
}

func NewCloser(l *zerolog.Logger, forClose ...C) *Closer {
	return &Closer{logger: l, forClose: forClose}
}

func (c *Closer) AppendClosers(forClose ...C) {
	c.forClose = append(c.forClose, forClose...)
}

func (c *Closer) AppendCloser(closer *Closer) {
	c.forClose = append(c.forClose, closer.forClose...)
}
