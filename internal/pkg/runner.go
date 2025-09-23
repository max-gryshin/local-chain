package pkg

import (
	"context"
	"log/slog"
	"sync"
)

type Runner interface {
	Run(context.Context) error
}

type Func func(context.Context) error

func (f Func) Run(ctx context.Context) error {
	return f(ctx)
}

func Run(parentCtx context.Context, logger *slog.Logger, runners ...Runner) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	wg := sync.WaitGroup{}

	errC := make(chan error, len(runners))

	for _, entry := range runners {
		entry := entry

		wg.Add(1)

		go func() {
			err := entry.Run(ctx)

			cancel()

			if err != nil {
				logger.Error(err.Error())
				errC <- err
			}

			wg.Done()
		}()
	}

	wg.Wait()

	var firstError error

	select {
	case firstError = <-errC:
	default:
	}

	return firstError
}
