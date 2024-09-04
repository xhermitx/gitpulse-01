package job

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/backend/types"
)

func TestJobHandlers(t *testing.T) {
	jobStore := &mockJobStore{}

	handler := NewHandler(jobStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8000/create", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/create", handler.CreateJobHandler)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Unexpected Status Code %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

}

type mockJobStore struct{}

func (s *mockJobStore) CreateJob(job types.Job) error {
	panic("uimplemented")
}

func (s *mockJobStore) UpdateJob(job types.Job) error {
	panic("uimplemented")
}

func (s *mockJobStore) DeleteJob(jobId uint) error {
	panic("uimplemented")
}

func (s *mockJobStore) ListJob() ([]types.Job, error) {
	panic("uimplemented")
}

func (s *mockJobStore) FindJobByName(name string) (*types.Job, error) {
	panic("unimplemented")
}
