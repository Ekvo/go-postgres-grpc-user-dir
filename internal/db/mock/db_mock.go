package mock

import (
	"context"
	"errors"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/model"
)

/*
type Provider interface {
	CreateUser(ctx context.Context, user *model.User) (uint, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUserByID(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	RemoveUserByID(ctx context.Context, id uint) error
	ClosePool()
}
*/

var ErrMockDB = errors.New("bad query")

type mockProvider struct {
	id          uint
	userByEmail map[string]*model.User
	userByID    map[uint]*model.User
}

func NewMockProvider() *mockProvider {
	return &mockProvider{
		userByEmail: make(map[string]*model.User),
		userByID:    make(map[uint]*model.User),
	}
}

func (mp *mockProvider) incrementID() {
	mp.id++
}

func (mp *mockProvider) CreateUser(_ context.Context, user *model.User) (uint, error) {
	if _, ex := mp.userByEmail[user.Email]; ex {
		return 0, ErrMockDB
	}
	mp.incrementID()
	user.ID = mp.id
	mp.userByEmail[user.Email] = user
	mp.userByID[user.ID] = user
	return user.ID, nil
}

func (mp *mockProvider) FindUserByEmail(_ context.Context, email string) (*model.User, error) {
	if user, ex := mp.userByEmail[email]; ex {
		return user, nil
	}
	return nil, ErrMockDB
}

func (mp *mockProvider) FindUserByID(_ context.Context, id uint) (*model.User, error) {
	if user, ex := mp.userByID[id]; ex {
		return user, nil
	}
	return nil, ErrMockDB
}

func (mp *mockProvider) UpdateUser(ctx context.Context, user *model.User) error {
	if _, ex := mp.userByID[user.ID]; ex {
		mp.userByID[user.ID] = user
		mp.userByEmail[user.Email] = user
		return nil
	}
	return ErrMockDB
}

func (mp *mockProvider) RemoveUserByID(ctx context.Context, id uint) error {
	if user, ex := mp.userByID[id]; ex {
		delete(mp.userByID, id)
		delete(mp.userByEmail, user.Email)
		return nil
	}
	return ErrMockDB
}

func (mp *mockProvider) ClosePool() {
}
