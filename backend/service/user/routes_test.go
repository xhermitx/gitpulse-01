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

func (s *mockUserStore) CreateUser(username string) (*types.User, error) {
	panic("unimplemented")
}

func (s *mockUserStore) DeleteUser(username string) error {
	panic("unimplemented")
}

func (s *mockUserStore) FindUserByEmail(email string) (*types.User, error) {
	panic("unimplemented")
}

func (s *mockUserStore) UpdateUserDetails(user types.User) error {
	panic("unimplemented")
}
