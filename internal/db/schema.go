// describes indexes and tables for a database
package db

import "context"

var (
	userTable = `CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(65) NOT NULL,
    first_name VARCHAR(128) NOT NULL,
    last_name VARCHAR(128) NULL,
    email VARCHAR(512) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL
);`

	loginBTreeIndex     = `CREATE INDEX IF NOT EXISTS login_btree_index ON users (login ASC);`
	emailBTreeIndex     = `CREATE INDEX IF NOT EXISTS email_btree_index ON users (email ASC);`
	createdAtBTreeIndex = `CREATE INDEX IF NOT EXISTS created_at_btree_index ON users (created_at DESC);`
)

func (p *provider) createTableIndex(ctx context.Context, schema ...string) error {
	for _, s := range schema {
		_, err := p.dbPool.Exec(ctx, s)
		if err != nil {
			return err
		}
	}
	return nil
}
