package models

// User is the base model
type User struct {
	Name     string `json:"name"`
	Cpf      string `json:"cpf"`
	Phone    string `json:"phone"`
	Priority bool   `json:"priority"`
}

// UserQueue is a queue of user
type UserQueue []User

// NewQueue returns a new UserQueue
func NewUserQueue() *UserQueue {
	return &UserQueue{}
}

// InsertNewUser insert a new user on the last position
func (q *UserQueue) InsertNewUser(user User) {
	*q = append(*q, user)
}

// RemoveFirstUser removes the first user on the top of the list
func (q *UserQueue) Remove() User {
	t := *q
	*q = t[1:]

	return t[0]
}
