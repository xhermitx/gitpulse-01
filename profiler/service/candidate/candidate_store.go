package candidate

import (
	"github.com/xhermitx/gitpulse-01/profiler/types"
	"gorm.io/gorm"
)

type CandidateStore struct {
	db *gorm.DB
}

func NewCandidateStore(db *gorm.DB) *CandidateStore {
	return &CandidateStore{
		db: db,
	}
}

func (s *CandidateStore) SaveCandidate(c types.Candidate) error {
	return nil
}
