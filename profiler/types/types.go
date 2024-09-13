package types

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CandidateStore interface {
	SaveCandidate(*Candidate) error
}

type Queue interface {
	Subscribe(queueName string) (<-chan amqp.Delivery, error)
}

type JobQueue struct {
	JobId     string   `json:"job_id"`
	JobDesc   string   `json:"job_desc"`
	Filename  string   `json:"filename"`
	GithubIDs []string `json:"github_ids"`
	Status    bool     `json:"status"`
}

type Cache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (value string, err error)
}

type GitService interface {
	FetchUserDetails(github_id string) (*GitUser, error)
}

type GitQuery struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

type GitResponse struct {
	Data struct {
		User GitUser `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type GitUser struct {
	CandidateMeta
	Followers          Followers        `json:"followers"`
	Contributions      ContributionData `json:"contributionsCollection"`
	TopRepo            RepositoryNode   `json:"repositories"`
	TopContributedRepo RepositoryNode   `json:"repositoriesContributedTo"`
}

type Followers struct {
	TotalCount int `json:"totalCount"`
}

type ContributionData struct {
	ContributionCalendar struct {
		TotalContributions int `json:"totalContributions"`
	} `json:"contributionCalendar"`
}

type RepositoryNode struct {
	Nodes []Repository `json:"nodes"`
}

type Repository struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	URL         string             `json:"url"`
	Stargazers  Stargazers         `json:"stargazers"`
	Languages   LanguageConnection `json:"languages"`
	Topics      RepositoryTopics   `json:"repositoryTopics"`
}

type Stargazers struct {
	TotalCount int `json:"totalCount"`
}

type LanguageConnection struct {
	Nodes []struct {
		Name string `json:"name"`
	} `json:"nodes"`
}

type RepositoryTopics struct {
	Nodes []struct {
		Topic struct {
			Name string `json:"name"`
		} `json:"topic"`
	} `json:"nodes"`
}

type CandidateMeta struct {
	Name        string `json:"name"`
	Username    string `json:"login"`
	AccountType string `json:"__typename" gorm:"-"` // This will hold "User" or "Organization"
	AvatarURL   string `json:"avatarUrl"`
	Bio         string `json:"bio"`
	// Company     string `json:"company" gorm:"-"`
	// Location    string `json:"location" gorm:"-"`
	Email      string `json:"email"`
	WebsiteURL string `json:"websiteUrl"`
}

type Candidate struct {
	CandidateMeta
	CandidateId             string `gorm:"primary_key"`
	TotalContributions      int
	TotalFollowers          int
	TopRepo                 string
	TopRepoStars            int
	TopContributedRepo      string
	TopContributedRepoStars int
	Languages               []string
	Topics                  []string
	JobId                   string `gorm:"unique"`
}
