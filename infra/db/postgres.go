package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/zer0day88/tinder/config"
	"github.com/zer0day88/tinder/migrations"
	"time"
)

func InitPostgres() (*pgxpool.Pool, error) {
	var (
		host              = config.Key.Database.Postgres.Host
		user              = config.Key.Database.Postgres.User
		password          = config.Key.Database.Postgres.Password
		port              = config.Key.Database.Postgres.Port
		dbname            = config.Key.Database.Postgres.DbName
		connMaxLifetime   = config.Key.Database.Postgres.ConMaxLifetime
		connMaxOpen       = config.Key.Database.Postgres.ConMaxOpen
		healthCheckPeriod = config.Key.Database.Postgres.HealthCheckPeriod
		connMaxIdleTime   = config.Key.Database.Postgres.ConMaxIdleTime
		sslMode           = config.Key.Database.Postgres.SSLMode
	)

	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		dbname,
		sslMode)

	cfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		panic(err)
	}

	cfg.MaxConns = connMaxOpen
	cfg.MaxConnLifetime = time.Duration(connMaxLifetime) * time.Second
	cfg.MaxConnIdleTime = time.Duration(connMaxIdleTime) * time.Second
	cfg.HealthCheckPeriod = time.Duration(healthCheckPeriod) * time.Second

	ctx := context.Background()

	dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping db failed")
	}

	return dbPool, nil
}

func MigratePG(postgres *pgxpool.Pool) {
	goose.SetBaseFS(migrations.EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	dba := stdlib.OpenDBFromPool(postgres)

	if err := goose.Up(dba, "."); err != nil {
		panic(err)
	}

	if err := dba.Close(); err != nil {
		panic(err)
	}
}
