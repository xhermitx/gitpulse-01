package user

import (
	"github.com/xhermitx/gitpulse-01/backend/types"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

// CreateUser implements types.UserStore.
func (s *Store) CreateUser(username string) (*types.User, error) {
	panic("unimplemented")
}

// DeleteUser implements types.UserStore.
func (s *Store) DeleteUser(username string) error {
	panic("unimplemented")
}

// FindUserByEmail implements types.UserStore.
func (s *Store) FindUserByEmail(email string) (*types.User, error) {
	panic("unimplemented")
}

// UpdateUserDetails implements types.UserStore.
func (s *Store) UpdateUserDetails(user types.User) error {
	panic("unimplemented")
}
