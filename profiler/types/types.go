package types

import amqp "github.com/rabbitmq/amqp091-go"

type CandidateStore interface {
	SaveCandidate(Candidate) error
}

type JobQueue struct {
	JobId     string
	Filename  string
	GithubIDs []string
	Status    bool
}

type Queue interface {
	Subscribe(queueName string) (<-chan amqp.Delivery, error)
}

type Cache interface {
}

type GitService interface {
	FetchUserDetails(github_id string) (*Candidate, error)
}

type GitQuery struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

type GitResponse struct {
	Data struct {
		Candidate Candidate `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type Candidate struct {
	Name                   string           `json:"name"`
	Username               string           `json:"login"`
	AccountType            string           `json:"__typename"` // This will hold "User" or "Organization"
	AvatarURL              string           `json:"avatarUrl"`
	Bio                    string           `json:"bio"`
	Company                string           `json:"company"`
	Location               string           `json:"location"`
	Email                  string           `json:"email"`
	WebsiteURL             string           `json:"websiteUrl"`
	Followers              Followers        `json:"followers"`
	Contributions          ContributionData `json:"contributionsCollection"`
	MostPopularRepo        RepositoryNode   `json:"repositories"`
	MostPopularContributed RepositoryNode   `json:"repositoriesContributedTo"`
}

type ContributionData struct {
	ContributionCalendar ContributionCalendar `json:"contributionCalendar"`
}

type ContributionCalendar struct {
	Weeks []Week `json:"weeks"`
}

type Week struct {
	ContributionDays []Contribution `json:"contributionDays"`
}

type Contribution struct {
	Date              string `json:"date"`
	ContributionCount int    `json:"contributionCount"`
}

type RepositoryNode struct {
	Nodes []Repository `json:"nodes"`
}

type Repository struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Stargazers  Stargazers         `json:"stargazers"`
	Languages   LanguageConnection `json:"languages"`
	Topics      RepositoryTopics   `json:"repositoryTopics"`
	URL         string             `json:"url"`
}

type Followers struct {
	TotalCount int `json:"totalCount"`
}

type Stargazers struct {
	TotalCount int `json:"totalCount"`
}

type LanguageConnection struct {
	Nodes []Language `json:"nodes"`
}

type Language struct {
	Name string `json:"name"`
}

type RepositoryTopics struct {
	Nodes []RepositoryTopicNode `json:"nodes"`
}

type RepositoryTopicNode struct {
	Topic Topic `json:"topic"`
}

type Topic struct {
	Name string `json:"name"`
}
