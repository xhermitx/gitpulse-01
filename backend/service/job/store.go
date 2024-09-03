package job

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

func (s *Store) CreateJob(job types.Job) (*types.Job, error) {
	panic("uimplemented")
}

func (s *Store) UpdateJob(job types.Job) (*types.Job, error) {
	panic("uimplemented")
}

func (s *Store) DeleteJob(id uint) error {
	panic("uimplemented")
}

func (s *Store) ListJob() ([]types.Job, error) {
	panic("uimplemented")
}
