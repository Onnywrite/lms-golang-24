package grace

import (
	"context"
	"io"

	"github.com/Onnywrite/lms-golang-24/pkg/erix"
)

type ShutdownGroup interface {
	io.Closer
	Add(io.Closer)
	WaitAndClose(context.Context) error
}

func NewShutdownGroup() ShutdownGroup {
	return &shutdownGroup{}
}

type shutdownGroup []io.Closer

func (s *shutdownGroup) Add(c io.Closer) {
	*s = append(*s, c)
}

func (s *shutdownGroup) Close() error {
	var errs []error

	for _, c := range *s {
		if err := c.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return erix.NewMulti(errs)
}

func (s *shutdownGroup) WaitAndClose(ctx context.Context) error {
	<-ctx.Done()

	err := s.Close()
	if err != nil {
		return err
	}

	return nil
}
