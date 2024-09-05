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

func (s *Store) CreateUser(user types.User) error {
	if res := s.db.Create(user); res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *Store) DeleteUser(userId string) error {
	if res := s.db.Delete(&types.User{}, userId); res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *Store) UpdateUser(user types.User) error {
	panic("unimplemented")
}

func (s *Store) FindUserById(userId string) (*types.User, error) {
	var user types.User
	if res := s.db.First(&user, userId); res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (s *Store) FindUserByEmail(email string) (*types.User, error) {
	var user types.User

	if res := s.db.First(&user, "email = ?", email); res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (s *Store) FindUserByUsername(username string) (*types.User, error) {
	var user types.User

	if res := s.db.First(&user, "username = ?", username); res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}
