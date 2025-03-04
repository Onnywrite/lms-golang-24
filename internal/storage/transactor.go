package storage

import "context"

type Transactor interface {
	Atomically(context.Context, func(ctx context.Context) error) error
}
