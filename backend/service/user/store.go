package user

import (
	"github.com/xhermitx/gitpulse-01/backend/types"
	"golang.org/x/crypto/bcrypt"
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
	panic("unimplemented")
}

func (s *Store) DeleteUser(userId string) error {
	panic("unimplemented")
}

func (s *Store) UpdateUser(user types.User) error {
	panic("unimplemented")
}

func (s *Store) LoginUser(credentials types.Credentials) (*types.User, error) {

	var user types.User
	switch true {
	case user.Email != "":
		if res := s.db.First(&user, "email = ?", credentials.Email); res.Error != nil {
			return nil, res.Error
		}

	case user.Username != "":
		if res := s.db.First(&user, "username = ?", credentials.Username); res.Error != nil {
			return nil, res.Error
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) FindUserById(userId string) (*types.User, error) {
	var user types.User
	if res := s.db.First(&user, userId); res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (s *Store) FindUserByEmail(email string) (*types.User, error) {
	panic("unimplemented")
}

func (s *Store) FindUserByUsername(username string) (*types.User, error) {
	panic("unimplemented")
}
