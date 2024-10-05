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
	UserId       string    `json:"user_id" gorm:"primary_key" example:"<user_id>"`
	FirstName    string    `json:"first_name" gorm:"not null" example:"John"`
	LastName     string    `json:"last_name" gorm:"not null" example:"Doe"`
	Username     string    `json:"username" gorm:"unique, not null" example:"jondo"`
	Email        string    `json:"email" gorm:"unique, not null" example:"doe.john@gmail.com"`
	Password     string    `json:"password" gorm:"not null" example:"johnkibilli@123"`
	Organization string    `json:"organization" gorm:"not null" example:"Illuminati"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime"`
}

type UserContext string

type Credentials struct {
	Email    string `json:"email" example:"doe.john@gmail.com"`
	Username string `json:"username" example:"jondo"`
	Password string `json:"password" example:"johnkibilli@123"`
}

type DeleteUserPayload struct {
	UserId string `json:"user_id" example:"<user_id>"`
}

type JobStore interface {
	CreateJob(Job) error
	UpdateJob(Job) error
	DeleteJob(string) error
	ListJobs(string) ([]Job, error)
	FindJobById(string, string) (*Job, error)
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
	JobId     string `json:"job_id" example:"<job_id>"`
	DriveLink string `json:"drive_link"`
}

type ParserPayload struct {
	JobId     string `json:"job_id"`
	DriveLink string `json:"drive_link"`
}

type DeleteJobPayload struct {
	JobId string `json:"job_id" example:"<job_id>"`
}

type CandidateStore interface {
	GetCandidateList(string) ([]results.Candidate, error)
}

type JobResultPayload struct {
	JobId string `json:"job_id"`
	Count int    `json:"count"`
}
