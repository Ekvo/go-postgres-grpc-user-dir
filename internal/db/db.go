package db

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/model"
)

type Provider interface {
	CreateUser(ctx context.Context, user *model.User) (uint, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUserByID(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	RemoveUserByID(ctx context.Context, id uint) error
	ClosePool()
}

type provider struct {
	dbPool *pgxpool.Pool
}

func OpenPool(ctx context.Context, cfg *config.Config) (*provider, error) {
	dbPool, err := initPool(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := dbPool.Ping(ctx); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("db: Ping error - %w", err)
	}
	provider := &provider{dbPool: dbPool}
	if err := provider.createTableIndex(
		ctx,
		userTable,
		loginBTreeIndex,
		emailBTreeIndex,
		createdAtBTreeIndex); err != nil {
		provider.ClosePool()
		return nil, fmt.Errorf("db: create schema error - %w", err)
	}
	return provider, nil
}

func (p *provider) ClosePool() {
	p.dbPool.Close()
	log.Print("db: *pgxpool.Pool is closed")
}

func initPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	cfgPgx, err := pgxpool.ParseConfig(cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("db: failed to parse pg config: %w", err)
	}
	cfgPgx.MaxConns = int32(10)
	cfgPgx.MinConns = int32(1)
	cfgPgx.HealthCheckPeriod = 1 * time.Minute
	cfgPgx.MaxConnLifetime = 24 * time.Hour
	cfgPgx.MaxConnIdleTime = 15 * time.Minute
	cfgPgx.ConnConfig.ConnectTimeout = 1 * time.Minute
	cfgPgx.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfgPgx.HealthCheckPeriod,
		Timeout:   cfgPgx.ConnConfig.ConnectTimeout,
	}).DialContext
	dbPool, err := pgxpool.NewWithConfig(ctx, cfgPgx)
	if err != nil {
		return nil, fmt.Errorf("db: failed to parse pg config: %w", err)
	}
	return dbPool, nil
}
