package store

import (
	"github.com/xhermitx/gitpulse-01/profiler/types"
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

func (s Store) SaveCandidate(c *types.Candidate) error {
	if res := s.db.Create(&c); res.Error != nil {
		return res.Error
	}
	return nil
}
