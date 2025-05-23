// describes rules for migration
package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/net/context"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
)

// Migration - contain connection to *sql.DB, *pgx.Config, driverName for db, dbURL, pathToMigrate
type Migration struct {
	conn   *sql.DB
	config *pgx.Config

	driverName string
	dbURL      string

	pathToMigrate string
}

// NewMigration - set data from config.MigrationConfig
func NewMigration(cfg *config.MigrationConfig) *Migration {
	return &Migration{
		config:        &pgx.Config{},
		driverName:    `postgres`,
		pathToMigrate: cfg.PathToMigrations,
		dbURL:         cfg.DBURL,
	}
}

// Up - member of Migration
// create connection to database
// check ping
// do migration
func (m *Migration) Up(ctx context.Context) error {
	if err := m.createConnect(); err != nil {
		return fmt.Errorf("migration: createConnect error - {%w};", err)
	}
	defer func() {
		_ = m.conn.Close()
		log.Printf("migration: conn is closed")
	}()

	if err := m.conn.PingContext(ctx); err != nil {
		return fmt.Errorf("migration: sql.DB.Ping error - {%w};", err)
	}
	log.Print("migration: ping is successful")

	if err := m.migrateSQL(); err != nil {
		return fmt.Errorf("migration: migrateSQL error - {%w};", err)
	}
	log.Printf("migration: migration is created")

	return nil
}

// migrateSQL is a member of migration
// use m.conn to create the driver (pgx.WithInstance)
// then create *migrate.Migrate and call mig.Up()
func (m *Migration) migrateSQL() error {
	driver, err := pgx.WithInstance(m.conn, m.config)
	if err != nil {
		return err
	}

	mig, err := migrate.NewWithDatabaseInstance(
		"file://"+m.pathToMigrate,
		m.driverName,
		driver,
	)
	if err != nil {
		return err
	}

	if err := mig.Up(); !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

// createConnect() - sql.Open with help m.driverName, m.dbURL and set m.conn
func (m *Migration) createConnect() error {
	conn, err := sql.Open(m.driverName, m.dbURL)
	if err != nil {
		return err
	}
	m.conn = conn

	log.Printf("migration: conn is created")

	return nil
}
