// initialization of database, with logic for query
package db

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/model"
)

// Provider - logic for work with store
type Provider interface {
	CreateUser(ctx context.Context, user *model.User) (uint, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUserByID(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	RemoveUserByID(ctx context.Context, id uint) error
	ClosePool()
}

// provider - wrapper for *pgxpool.Pool
type provider struct {
	dbPool *pgxpool.Pool
}

// OpenPool - call initPool to open pgx.pool, check Ping, create tables, indexes for the database if they do not exist
func OpenPool(ctx context.Context, cfg *config.DataBaseConfig) (*provider, error) {
	dbPool, err := initPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := dbPool.Ping(ctx); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("db: Ping error - %w", err)
	}
	log.Print("db: ping is successful")

	provider := &provider{dbPool: dbPool}

	return provider, nil
}

func (p *provider) ClosePool() {
	p.dbPool.Close()
	log.Print("db: database is closed")
}

// initPool - parse DBURL, set pgxpool.Config, open pgx.pool
func initPool(ctx context.Context, cfg *config.DataBaseConfig) (*pgxpool.Pool, error) {
	cfgPgx, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("db: failed to parse pg config: %w", err)
	}

	cfgPgx.MaxConns = int32(cfg.MaxConn)
	cfgPgx.MinConns = int32(cfg.MaxConn)
	cfgPgx.HealthCheckPeriod = cfg.HealthCheckPeriod
	cfgPgx.MaxConnLifetime = cfg.ConnMaxLifeTime
	cfgPgx.MaxConnIdleTime = cfg.ConnMaxIdleTime
	cfgPgx.ConnConfig.ConnectTimeout = cfg.ConnTime
	cfgPgx.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfgPgx.HealthCheckPeriod,
		Timeout:   cfgPgx.ConnConfig.ConnectTimeout,
	}).DialContext
	dbPool, err := pgxpool.NewWithConfig(ctx, cfgPgx)
	if err != nil {
		return nil, fmt.Errorf("db: failed to parse pg config: %w", err)
	}

	log.Print("db: database connected")

	return dbPool, nil
}
