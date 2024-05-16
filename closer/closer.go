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

type Closer struct {
	logger   *zerolog.Logger
	forClose []io.Closer
}

func (c *Closer) Close() {
	for _, closer := range c.forClose {
		err := closer.Close()
		if err != nil {
			c.logger.Err(err).Msg("failed to close:")
		}
	}

	c.logger.Info().Msg("all dependencies are closed")
}

func NewCloser(l *zerolog.Logger, forClose ...io.Closer) *Closer {
	return &Closer{logger: l, forClose: forClose}
}

func (c *Closer) AppendClosers(forClose ...io.Closer) {
	c.forClose = append(c.forClose, forClose...)
}

func (c *Closer) AppendCloser(closer *Closer) {
	c.forClose = append(c.forClose, closer.forClose...)
}
