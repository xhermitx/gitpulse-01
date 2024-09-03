package types

type UserStore interface {
	CreateUser(username string) (*User, error)
	DeleteUser(username string) error
	FindUserByEmail(email string) (*User, error)
	UpdateUserDetails(user User) error
}

type User struct {
	UserID       uint
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Organization string `json:"organization"`
	CreatedAt    string
}

type JobStore interface {
	CreateJob(job Job) (*Job, error)
	UpdateJob(job Job) (*Job, error)
	DeleteJob(id uint) error
	ListJob() ([]Job, error)
}

type Job struct {
	JobID       uint
	Name        string `json:"name"`
	Description string `json:"description"`
	DriveLink   string `json:"drive_link"`
	CreatedAt   string
}

type CandidateList struct {
	CandidateID    uint   `json:"candidate_id"`
	GithubUsername string `json:"github_username"`
	Score          int    `json:"score"`
}
