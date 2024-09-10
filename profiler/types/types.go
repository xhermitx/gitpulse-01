package types

import amqp "github.com/rabbitmq/amqp091-go"

type CandidateStore interface {
	SaveCandidate(Candidate) error
}

type Candidate struct {
	CandidateId     string `json:"candidate_id" gorm:"primary_key"`
	GithubId        string `json:"github_id" gorm:"unique, not null"`
	Followers       uint   `json:"followers" gorm:"not null"`
	Contributions   uint   `json:"contributions" gorm:"not null"`
	MostPopularRepo string `json:"most_popular_repo" gorm:"not null"`
	RepoStars       uint   `json:"repo_stars" gorm:"not null"`
	Score           int    `json:"score" gorm:"not null"`
	JobId           string `json:"job_id" gorm:"not null"`
}

type GitService interface {
	FetchUserDetails(github_id string) (Candidate, error)
}

type Queue interface {
	Subscribe(queueName string) (<-chan amqp.Delivery, error)
}

type Cache interface {
}
