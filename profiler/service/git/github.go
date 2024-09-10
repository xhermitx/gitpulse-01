package git

import "github.com/xhermitx/gitpulse-01/profiler/types"

type GitService struct {
}

func NewGitService() *GitService {
	return &GitService{}
}

func (g *GitService) FetchUserDetails(github_id string) (types.Candidate, error) {
	return types.Candidate{}, nil
}
