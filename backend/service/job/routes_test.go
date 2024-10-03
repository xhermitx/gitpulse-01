package job

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/backend/types"
	results "github.com/xhermitx/gitpulse-results"
)

func TestJobHandlers(t *testing.T) {
	jobStore := &mockJobStore{}
	userStore := &mockUserStore{}
	candidateStore := &mockCandidateStore{}

	handler := NewHandler(jobStore, userStore, candidateStore)

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

func (s *mockJobStore) CreateJob(_ types.Job) error {
	panic("uimplemented")
}
func (s *mockJobStore) UpdateJob(_ types.Job) error {
	panic("uimplemented")
}
func (s *mockJobStore) DeleteJob(_ string) error {
	panic("uimplemented")
}
func (s *mockJobStore) ListJobs(_ string) ([]types.Job, error) {
	panic("uimplemented")
}
func (s *mockJobStore) FindJobById(_, _ string) (*types.Job, error) {
	panic("unimplemented")
}

type mockUserStore struct{}

func (s *mockUserStore) CreateUser(_ types.User) error {
	panic("unimplemented")
}
func (s *mockUserStore) DeleteUser(_ string) error {
	panic("unimplemented")
}
func (s *mockUserStore) UpdateUser(_ types.User) error {
	panic("unimplemented")
}
func (s *mockUserStore) FindUserByEmail(_ string) (*types.User, error) {
	panic("unimplemented")
}
func (s *mockUserStore) FindUserById(_ string) (*types.User, error) {
	panic("unimplemented")
}
func (s *mockUserStore) FindUserByUsername(_ string) (*types.User, error) {
	panic("unimplemented")
}

type mockCandidateStore struct{}

func (s mockCandidateStore) GetCandidateList(_ string) ([]results.Candidate, error) {
	panic("unimplemented")
}
