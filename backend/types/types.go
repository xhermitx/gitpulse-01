package types

import "time"

type UserStore interface {
	CreateUser(User) error
	DeleteUser(string) error
	UpdateUser(User) error

	FindUserById(string) (*User, error)
	FindUserByEmail(string) (*User, error)
	FindUserByUsername(string) (*User, error)
}

type User struct {
	UserId       string    `json:"user_id,omitempty" gorm:"primary_key"`
	FirstName    string    `json:"first_name,omitempty" gorm:"not null"`
	LastName     string    `json:"last_name,omitempty" gorm:"not null"`
	Username     string    `json:"username,omitempty" gorm:"unique, not null"`
	Email        string    `json:"email,omitempty" gorm:"unique, not null"`
	Password     string    `json:"password,omitempty" gorm:"not null"`
	Organization string    `json:"organization,omitempty" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at,omitempty" gorm:"type:datetime"`
}

type UserContext string

type Credentials struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type JobStore interface {
	CreateJob(Job) error
	UpdateJob(Job) error
	DeleteJob(string) error
	ListJobs(string) ([]Job, error)
	FindJobById(string) (*Job, error)
}

type Job struct {
	JobId       string    `json:"job_id,omitempty" gorm:"primary_key"`
	JobName     string    `json:"job_name,omitempty" gorm:"unique, not null"`
	Description string    `json:"description,omitempty" gorm:"not null"`
	DriveLink   string    `json:"drive_link,omitempty" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at,omitempty" gorm:"type:datetime"`
	UserId      string    `json:"user_id" gorm:"not null"`
}

type DeleteJobPayload struct {
	JobId string `json:"job_id"`
}

type DeleteUserPayload struct {
	UserId string `json:"user_id"`
}

type CandidateStore interface {
	GetCandidateList(string) ([]*Candidate, error)
}

type CandidateMeta struct {
	Name       string `json:"name"`
	Username   string `json:"login"`
	AvatarURL  string `json:"avatarUrl"`
	Bio        string `json:"bio"`
	Email      string `json:"email"`
	WebsiteURL string `json:"websiteUrl"`
}

type Candidate struct {
	CandidateMeta
	CandidateId             string   `json:"candidate_id" gorm:"primary_key"`
	TotalContributions      int      `json:"total_contributions"`
	TotalFollowers          int      `json:"total_followers"`
	TopRepo                 string   `json:"top_repo"`
	TopRepoStars            int      `json:"top_repo_stars"`
	TopContributedRepo      string   `json:"top_contributed_repo"`
	TopContributedRepoStars int      `json:"top_contributed_repo_stars"`
	Languages               []string `json:"languages"`
	Topics                  []string `json:"topics"`
	JobId                   string   `json:"job_id" gorm:"unique"`
	Score                   int      `json:"score" gorm:"-"`
}

type JobResultPayload struct {
	JobId string `json:"job_id"`
	Count string `json:"count"`
}
