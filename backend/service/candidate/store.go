package candidate

import (
	"github.com/xhermitx/gitpulse-01/backend/types"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return Store{
		db: db,
	}
}

func (s Store) GetCandidateList(jobId string) ([]*types.Candidate, error) {
	var list []*types.Candidate
	if res := s.db.Find(list, "job_id = ?", jobId); res.Error != nil {
		return nil, res.Error
	}
	return list, nil
}
