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

func (s *Store) CreateJob(job types.Job) error {
	if res := s.db.Create(&job); res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *Store) UpdateJob(job types.Job) error {
	res := s.db.Save(&job)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *Store) DeleteJob(jobId string) error {
	if res := s.db.Delete(&types.Job{}, "job_id", jobId); res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *Store) ListJobs(userId string) ([]types.Job, error) {
	var job []types.Job
	if res := s.db.Find(&job, "user_id = ?", userId); res.Error != nil {
		return nil, res.Error
	}
	return job, nil
}

func (s *Store) FindJobById(jobId string, userId string) (*types.Job, error) {
	var job types.Job
	if res := s.db.First(&job, "job_id = ? AND user_id = ?", jobId, userId); res.Error != nil {
		return nil, res.Error
	}
	return &job, nil
}
