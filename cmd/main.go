package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/zer0day88/tinder/config"
	"github.com/zer0day88/tinder/infra/db"
	"github.com/zer0day88/tinder/internal/app/route"
	"github.com/zer0day88/tinder/pkg/logger"
	"github.com/zer0day88/tinder/pkg/shutdown"
	"net/http"
	"time"
)

func main() {

	config.Load()
	log := logger.New()

	postgres, err := db.InitPostgres()
	if err != nil {
		panic(err)
	}

	db.MigratePG(postgres)

	rdb, err := db.InitRedis()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	route.InitRoute(e, postgres, log, rdb)

	// Start server
	go func() {
		if err := e.Start(":" + config.Key.ServerPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	wait := shutdown.GracefulShutdown(context.Background(), log, 10*time.Second, map[string]shutdown.Operation{

		"http-server": func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
		"database": func(ctx context.Context) error {
			postgres.Close()
			return nil
		},
		// Add other cleanup operations here
	})

	<-wait
}
