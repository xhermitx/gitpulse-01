package types

type UserStore interface {
	CreateUser(username string) (*User, error)
	DeleteUser(username string) error
	FindUserByEmail(email string) (*User, error)
	UpdateUserDetails(user User) error
}
type User struct {
	UserId       string `json:"user_id,omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	FirstName    string `json:"first_name,omitempty" gorm:"not null"`
	LastName     string `json:"last_name,omitempty" gorm:"not null"`
	Username     string `json:"username,omitempty" gorm:"unique, not null"`
	Email        string `json:"email,omitempty" gorm:"unique, not null"`
	Password     string `json:"password,omitempty" gorm:"not null"`
	Organization string `json:"organization,omitempty" gorm:"not null"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type JobStore interface {
	CreateJob(Job) error
	UpdateJob(Job) error
	DeleteJob(uint) error
	ListJob() ([]Job, error)
	FindJobByName(string) (*Job, error)
}

type Job struct {
	JobId       string `json:"job_id,omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	Name        string `json:"name,omitempty" gorm:"unique, not null"`
	Description string `json:"description,omitempty" gorm:"not null"`
	DriveLink   string `json:"drive_link,omitempty" gorm:"not null"`
	UserId      string `json:"omitempty"`
}

type Candidate struct {
	CandidateId     uint   `json:"candidate_id,omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	GithubId        string `json:"github_username,omitempty" gorm:"unique, not null"`
	Followers       uint   `json:"followers,omitempty" gorm:"not null"`
	Contributions   uint   `json:"contributions,omitempty" gorm:"not null"`
	MostPopularRepo string `json:"most_popular_repo,omitempty" gorm:"not null"`
	RepoStars       uint   `json:"repo_stars,omitempty" gorm:"not null"`
	Score           int    `json:"score,omitempty" gorm:"not null"`
	JobId           string `json:"job_id,omitempty"`
}
