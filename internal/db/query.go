package db

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/model"
)

func (p *provider) CreateUser(ctx context.Context, user *model.User) (uint, error) {
	userID := uint(0)
	err := p.dbPool.QueryRow(ctx, `
INSERT INTO users (
                   login,
                   password,
                   first_name,
                   last_name,
                   email,
                   created_at
                   )
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING id;`,
		user.Login,                             //1
		user.Password,                          //2
		user.FirstName,                         //3
		whenStringEmptyThenNULL(user.LastName), //4
		user.Email,                             //5
		user.CreatedAt,                         //6
	).Scan(&userID)
	return userID, err
}

func (p *provider) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	row := p.dbPool.QueryRow(ctx, `
SELECT * 
FROM users
WHERE email = $1
LIMIT 1;`, email)
	return scanUser(row)
}

func (p *provider) FindUserByID(ctx context.Context, id uint) (*model.User, error) {
	row := p.dbPool.QueryRow(ctx, `
SELECT * 
FROM users
WHERE id = $1
LIMIT 1;`, id)
	return scanUser(row)
}

func (p *provider) UpdateUser(ctx context.Context, user *model.User) error {
	upID := uint(0)
	err := p.dbPool.QueryRow(ctx, `
UPDATE users
SET login = $2,
    password = $3,
    first_name = $4,
    last_name = $5,
    email = $6,
    updated_at = $7
WHERE id = $1
RETURNING id;`,
		user.ID,                                //1
		user.Login,                             //2
		user.Password,                          //3
		user.FirstName,                         //4
		whenStringEmptyThenNULL(user.LastName), //5
		user.Email,                             //6
		user.UpdatedAt,                         //7
	).Scan(&upID)
	return err
}
func (p *provider) RemoveUserByID(ctx context.Context, id uint) error {
	delID := uint(0)
	err := p.dbPool.QueryRow(ctx, `
DELETE
FROM users
WHERE id = $1
RETURNING id;`, id).Scan(&delID)
	return err
}

func scanUser(row pgx.Row) (*model.User, error) {
	var (
		user model.User

		lastName  sql.NullString
		updatedAt sql.NullTime
	)
	if err := row.Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.FirstName,
		&lastName,
		&user.Email,
		&user.CreatedAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	if lastName.Valid {
		user.LastName = lastName.String
	}
	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}
	return &user, nil
}

func whenStringEmptyThenNULL(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
