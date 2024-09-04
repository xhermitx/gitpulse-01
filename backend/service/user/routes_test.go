package user

import (
	"testing"

	"github.com/xhermitx/gitpulse-01/backend/types"
)

func TestUserServiceHandlers(t *testing.T) {

	userStore := &mockUserStore{}

	handler := NewHandler(userStore)

	_ = handler
}

type mockUserStore struct{}

func (s *mockUserStore) CreateUser(_ types.User) error {
	panic("unimplemented")
}

func (s *mockUserStore) DeleteUser(_ string) error {
	panic("unimplemented")
}

func (s *mockUserStore) UpdateUser(_ types.User) error {
	panic("unimplemented")
}

func (s *mockUserStore) FindUserByEmail(_ string) (*types.User, error) {
	panic("unimplemented")
}

func (s *mockUserStore) FindUserById(_ string) (*types.User, error) {
	panic("unimplemented")
}

func (s *mockUserStore) FindUserByUsername(_ string) (*types.User, error) {
	panic("unimplemented")
}

func (s *mockUserStore) LoginUser(credentials types.Credentials) (*types.User, error) {
	panic("unimplemented")
}
