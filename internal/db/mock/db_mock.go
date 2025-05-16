package mock

import (
	"context"
	"errors"

	"github.com/Ekvo/go-postgres-grpc-user-dir/internal/model"
)

// ErrMockDB - imitation error from mockbase
var ErrMockDB = errors.New("bad query")

type mockProvider struct {
	id          uint
	userByID    map[uint]*model.User
	userByEmail map[string]*model.User
	userLogin   map[string]*model.User
}

func NewMockProvider() *mockProvider {
	return &mockProvider{
		userByID:    make(map[uint]*model.User),
		userByEmail: make(map[string]*model.User),
		userLogin:   make(map[string]*model.User),
	}
}

func (mp *mockProvider) incrementID() {
	mp.id++
}

func (mp *mockProvider) CreateUser(_ context.Context, user *model.User) (uint, error) {
	if _, ex := mp.userLogin[user.Login]; ex {
		return 0, ErrMockDB
	}
	if _, ex := mp.userByEmail[user.Email]; ex {
		return 0, ErrMockDB
	}
	mp.incrementID()
	user.ID = mp.id
	mp.createUser(user)
	return user.ID, nil
}

func (mp *mockProvider) createUser(user *model.User) {
	mp.userByID[user.ID] = user
	mp.userByEmail[user.Email] = user
	mp.userLogin[user.Login] = user
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
	if userEmail, ex := mp.userByEmail[user.Email]; ex && userEmail.ID != user.ID {
		return ErrMockDB
	}
	if userLogin, ex := mp.userLogin[user.Login]; ex && userLogin.ID != user.ID {
		return ErrMockDB
	}
	if _, ex := mp.userByID[user.ID]; ex {
		mp.userByID[user.ID] = user
		mp.userByEmail[user.Email] = user
		mp.userLogin[user.Login] = user
		return nil
	}
	return ErrMockDB
}

func (mp *mockProvider) RemoveUserByID(ctx context.Context, id uint) error {
	if user, ex := mp.userByID[id]; ex {
		delete(mp.userByID, id)
		delete(mp.userByEmail, user.Email)
		delete(mp.userLogin, user.Login)
		return nil
	}
	return ErrMockDB
}

func (mp *mockProvider) ClosePool() {
}
