package job

import (
	"testing"

	"github.com/xhermitx/gitpulse-01/backend/types"
)

func TestJobHandlers(t *testing.T) {
	jobStore := &mockJobStore{}

	handler := NewHandler(jobStore)

	_ = handler
}

type mockJobStore struct{}

func (s *mockJobStore) CreateJob(job types.Job) (*types.Job, error) {
	panic("uimplemented")
}

func (s *mockJobStore) UpdateJob(job types.Job) (*types.Job, error) {
	panic("uimplemented")
}

func (s *mockJobStore) DeleteJob(id uint) error {
	panic("uimplemented")
}

func (s *mockJobStore) ListJob() ([]types.Job, error) {
	panic("uimplemented")
}
