package types

import (
	"time"

	results "github.com/xhermitx/gitpulse-results"
)

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

type TriggerJobPayload struct {
	JobId     string `json:"job_id"`
	DriveLink string `json:"drive_link"`
}

type ParserPayload struct {
	JobId     string `json:"job_id"`
	DriveLink string `json:"drive_link"`
}

type DeleteJobPayload struct {
	JobId string `json:"job_id"`
}

type DeleteUserPayload struct {
	UserId string `json:"user_id"`
}

type CandidateStore interface {
	GetCandidateList(string) ([]results.Candidate, error)
}

type JobResultPayload struct {
	JobId string `json:"job_id"`
	Count int    `json:"count"`
}
