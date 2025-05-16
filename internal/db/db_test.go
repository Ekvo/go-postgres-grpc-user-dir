package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/config"
	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/model"
)

func newProviderForTest(ctx context.Context) (*provider, error) {
	pr, err := OpenPool(ctx, &config.Config{
		DBURL: `postgresql://manager:qwert12345@127.0.0.1:5432/testdb`,
	})
	if err != nil {
		return nil, err
	}

	_, err = pr.dbPool.Exec(ctx, `DELETE FROM users;`)

	return pr, err
}

func TestProvider_CreateUser(t *testing.T) {
	log.Printf("db_test: TestProvider_CreateUser - START")

	asserts := assert.New(t)
	requires := require.New(t)

	var testData = []struct {
		title       string
		logicOfTest func(ctx context.Context, pr *provider) (uint, error)
		expectedRes string
		err         error
		msg         string
	}{
		{
			title: `valid create user`,
			logicOfTest: func(ctx context.Context, pr *provider) (uint, error) {
				user := &model.User{
					Login:     `alien`,
					Password:  `avp`,
					FirstName: `Alex`,
					Email:     `alex@example.com`,
					CreatedAt: time.Now(),
				}
				return pr.CreateUser(ctx, user)
			},
			expectedRes: `^[1-9][0-9]*$`,
			err:         nil,
			msg:         `this must be valid, return uint > 0, error is nul`,
		},
		{
			title: `valid create user (all fields)`,
			logicOfTest: func(ctx context.Context, pr *provider) (uint, error) {
				create := time.Now()
				update := create.Add(time.Hour)
				user := &model.User{
					Login:     `foxy`,
					Password:  `qwert12345`,
					FirstName: `Alex`,
					LastName:  `Doe`,
					Email:     `doeFox@example.com`,
					CreatedAt: create,
					UpdatedAt: &update,
				}
				return pr.CreateUser(ctx, user)
			},
			expectedRes: `^[1-9][0-9]*$`,
			err:         nil,
			msg:         `this must be valid, return uint > 0, error is nul`,
		},
		{
			title: `wrong create, login already exist`,
			logicOfTest: func(ctx context.Context, pr *provider) (uint, error) {
				user := &model.User{
					Login:     `approve1`,
					Password:  `qwer`,
					FirstName: `Tom1`,
					Email:     `Tom1@example.com`,
					CreatedAt: time.Now(),
				}
				if _, err := pr.CreateUser(ctx, user); err != nil {
					return 0, fmt.Errorf("test ended early - %w", err)
				}
				user.Email = `Tom@example.com`
				return pr.CreateUser(ctx, user)
			},
			expectedRes: `0`,
			err:         errors.New("duplicate key value violates unique constraint \"users_login_key\""),
			msg:         `wrong, return uint = 0, error is exist`,
		},
		{
			title: `wrong create, email already exist`,
			logicOfTest: func(ctx context.Context, pr *provider) (uint, error) {
				user := &model.User{
					Login:     `approve2`,
					Password:  `qwer`,
					FirstName: `Tom2`,
					Email:     `Tom2@example.com`,
					CreatedAt: time.Now(),
				}
				if _, err := pr.CreateUser(ctx, user); err != nil {
					return 0, fmt.Errorf("test ended early - %w", err)
				}
				user.Login = `no approve`
				return pr.CreateUser(ctx, user)
			},
			expectedRes: `0`,
			err:         errors.New("duplicate key value violates unique constraint \"users_email_key\""),
			msg:         `wrong, return uint = 0, error is exist`,
		},
	}

	ctx := context.Background()

	pr, err := newProviderForTest(ctx)
	requires.NoError(err, "wrong connect to db")
	defer pr.ClosePool()

	for i, test := range testData {
		log.Printf("\t%d %s", i+1, test.title)

		res, err := test.logicOfTest(ctx, pr)

		if err == nil {
			asserts.ErrorIs(err, test.err, test.msg)
		} else {
			requires.Error(test.err, test.msg)
			asserts.Regexp(test.err.Error(), err.Error(), test.msg)
		}
		asserts.Regexp(test.expectedRes, res, test.msg)
	}
	log.Printf("db_test: TestProvider_CreateUser - END")
}

func TestProvider_FindUserByEmail(t *testing.T) {
	log.Printf("db_test: TestProvider_FindUserByEmail - START")

	asserts := assert.New(t)
	requires := require.New(t)

	var testData = []struct {
		title       string
		logicOfTest func(ctx context.Context, pr *provider) (*model.User, error)
		expectedRes *model.User
		err         error
		msg         string
	}{
		{
			title: `valid find user`,
			logicOfTest: func(ctx context.Context, pr *provider) (*model.User, error) {
				user := &model.User{
					Login:     `alien`,
					Password:  `avp`,
					FirstName: `Alex`,
					Email:     `alex@example.com`,
					CreatedAt: time.Now(),
				}
				if _, err := pr.CreateUser(ctx, user); err != nil {
					return nil, err
				}
				return pr.FindUserByEmail(ctx, user.Email)
			},
			expectedRes: &model.User{
				Login:     `alien`,
				Password:  `avp`,
				FirstName: `Alex`,
				Email:     `alex@example.com`,
			},
			err: nil,
			msg: `this must be valid, return ptrUser, error is nul`,
		},
		{
			title: `wrong find not exist`,
			logicOfTest: func(ctx context.Context, pr *provider) (*model.User, error) {
				return pr.FindUserByEmail(ctx, `alien@example.com`)
			},
			expectedRes: nil,
			err:         pgx.ErrNoRows,
			msg:         `wrong find by email not found, error is exist`,
		},
	}

	ctx := context.Background()

	pr, err := newProviderForTest(ctx)
	requires.NoError(err, "wrong connect to db")
	defer pr.ClosePool()

	for i, test := range testData {
		log.Printf("\t%d %s", i+1, test.title)

		res, err := test.logicOfTest(ctx, pr)

		if err == nil {
			asserts.ErrorIs(err, test.err, test.msg)
		} else {
			requires.Error(test.err, test.msg)
			asserts.Regexp(test.err.Error(), err.Error(), test.msg)
		}

		if res == nil {
			requires.Nil(test.expectedRes, test.msg)
		} else {
			test.expectedRes.ID = res.ID
			test.expectedRes.CreatedAt = res.CreatedAt

			asserts.Equal(test.expectedRes, res, test.msg)
		}

	}
	log.Printf("db_test: TestProvider_FindUserByEmail - END")
}

func TestProvider_FindUserByID(t *testing.T) {
	log.Printf("db_test: TestProvider_FindUserByID - START")

	asserts := assert.New(t)
	requires := require.New(t)

	var testData = []struct {
		title       string
		logicOfTest func(ctx context.Context, pr *provider) (*model.User, error)
		expectedRes *model.User
		err         error
		msg         string
	}{
		{
			title: `valid find user`,
			logicOfTest: func(ctx context.Context, pr *provider) (*model.User, error) {
				create := time.Now()
				update := create.Add(time.Hour)
				user := &model.User{
					Login:     `alien`,
					Password:  `avp`,
					FirstName: `Alex`,
					LastName:  `Vense`,
					Email:     `alex@example.com`,
					CreatedAt: create,
				}
				id, err := pr.CreateUser(ctx, user)
				if err != nil {
					return nil, err
				}
				user.ID = id
				user.UpdatedAt = &update
				if err := pr.UpdateUser(ctx, user); err != nil {
					return nil, err
				}
				return pr.FindUserByID(ctx, id)
			},
			expectedRes: &model.User{
				Login:     `alien`,
				Password:  `avp`,
				FirstName: `Alex`,
				LastName:  `Vense`,
				Email:     `alex@example.com`,
			},
			err: nil,
			msg: `this must be valid, return ptrUser, error is nul`,
		},
		{
			title: `wrong find not exist`,
			logicOfTest: func(ctx context.Context, pr *provider) (*model.User, error) {
				return pr.FindUserByID(ctx, 1)
			},
			expectedRes: nil,
			err:         pgx.ErrNoRows,
			msg:         `wrong find by ID not found, error is exist`,
		},
	}

	ctx := context.Background()

	pr, err := newProviderForTest(ctx)
	requires.NoError(err, "wrong connect to db")
	defer pr.ClosePool()

	for i, test := range testData {
		log.Printf("\t%d %s", i+1, test.title)

		res, err := test.logicOfTest(ctx, pr)

		if err == nil {
			asserts.ErrorIs(err, test.err, test.msg)
		} else {
			requires.Error(test.err, test.msg)
			asserts.Regexp(test.err.Error(), err.Error(), test.msg)
		}

		if res == nil {
			requires.Nil(test.expectedRes, test.msg)
		} else {
			test.expectedRes.ID = res.ID
			test.expectedRes.CreatedAt = res.CreatedAt
			test.expectedRes.UpdatedAt = res.UpdatedAt

			asserts.Equal(test.expectedRes, res, test.msg)
		}

	}
	log.Printf("db_test: TestProvider_FindUserByID - END")
}

func TestProvider_UpdateUser(t *testing.T) {
	log.Printf("db_test: TestProvider_UpdateUser - START")

	asserts := assert.New(t)
	requires := require.New(t)

	var testData = []struct {
		title       string
		logicOfTest func(ctx context.Context, pr *provider) error
		err         error
		msg         string
	}{
		{
			title: `valid update user`,
			logicOfTest: func(ctx context.Context, pr *provider) error {
				create := time.Now()
				update := create.Add(time.Hour)
				user := &model.User{
					Login:     `alien`,
					Password:  `avp`,
					FirstName: `Alex`,
					Email:     `alex@example.com`,
					CreatedAt: create,
				}
				id, err := pr.CreateUser(ctx, user)
				if err != nil {
					return err
				}
				user.ID = id
				user.LastName = `predator`
				user.UpdatedAt = &update
				return pr.UpdateUser(ctx, user)
			},
			err: nil,
			msg: `update must be valid, error is nul`,
		},
		{
			title: `invalid update, user not exist`,
			logicOfTest: func(ctx context.Context, pr *provider) error {
				create := time.Now()
				update := create.Add(time.Hour)
				user := &model.User{
					ID:        1,
					Login:     `alien`,
					Password:  `avp`,
					FirstName: `Alex`,
					Email:     `alex@example.com`,
					CreatedAt: create,
					UpdatedAt: &update,
				}
				return pr.UpdateUser(ctx, user)
			},
			err: pgx.ErrNoRows,
			msg: `wrong update , error is exist`,
		},
	}

	ctx := context.Background()

	pr, err := newProviderForTest(ctx)
	requires.NoError(err, "wrong connect to db")
	defer pr.ClosePool()

	for i, test := range testData {
		log.Printf("\t%d %s", i+1, test.title)

		err := test.logicOfTest(ctx, pr)

		if err == nil {
			asserts.ErrorIs(err, test.err, test.msg)
		} else {
			requires.Error(test.err, test.msg)
			asserts.Regexp(test.err.Error(), err.Error(), test.msg)
		}
	}
	log.Printf("db_test: TestProvider_UpdateUser - END")
}

func TestProvider_DeleteUser(t *testing.T) {
	log.Printf("db_test: TestProvider_DeleteUser - START")

	asserts := assert.New(t)
	requires := require.New(t)

	var testData = []struct {
		title       string
		logicOfTest func(ctx context.Context, pr *provider) error
		err         error
		msg         string
	}{
		{
			title: `valid user delete`,
			logicOfTest: func(ctx context.Context, pr *provider) error {
				user := &model.User{
					Login:     `alien`,
					Password:  `avp`,
					FirstName: `Alex`,
					Email:     `alex@example.com`,
					CreatedAt: time.Now(),
				}
				id, err := pr.CreateUser(ctx, user)
				if err != nil {
					return err
				}
				return pr.RemoveUserByID(ctx, id)
			},
			err: nil,
			msg: `update must be valid, error is nul`,
		},
		{
			title: `invalid delete, user not exist`,
			logicOfTest: func(ctx context.Context, pr *provider) error {
				return pr.RemoveUserByID(ctx, 1)
			},
			err: pgx.ErrNoRows,
			msg: `wrong delete, error is exist`,
		},
	}

	ctx := context.Background()

	pr, err := newProviderForTest(ctx)
	requires.NoError(err, "wrong connect to db")
	defer pr.ClosePool()

	for i, test := range testData {
		log.Printf("\t%d %s", i+1, test.title)

		err := test.logicOfTest(ctx, pr)

		if err == nil {
			asserts.ErrorIs(err, test.err, test.msg)
		} else {
			requires.Error(test.err, test.msg)
			asserts.Regexp(test.err.Error(), err.Error(), test.msg)
		}
	}
	log.Printf("db_test: TestProvider_DeleteUser - END")
}
