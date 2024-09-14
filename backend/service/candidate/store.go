package candidate

import (
	results "github.com/xhermitx/gitpulse-results"
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

func (s Store) GetCandidateList(jobId string) ([]*results.Candidate, error) {
	var list []*results.Candidate
	if res := s.db.Find(list, "job_id = ?", jobId); res.Error != nil {
		return nil, res.Error
	}
	return list, nil
}
