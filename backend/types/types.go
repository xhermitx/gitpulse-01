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
