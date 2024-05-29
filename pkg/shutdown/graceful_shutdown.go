package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type Operation func(ctx context.Context) error

func GracefulShutdown(ctx context.Context, log zerolog.Logger, timeout time.Duration, ops map[string]Operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Warn().Msg("shutting down")

		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Error().Msgf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Warn().Msgf("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Error().Err(err).Msgf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Warn().Msgf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
