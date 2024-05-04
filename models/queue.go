package models

import "time"

// User is the base model
type User struct {
	Name     string    `json:"name"`
	Cpf      string    `json:"cpf"`
	Phone    string    `json:"phone"`
	Priority bool      `json:"priority"`
	EnterOn  time.Time `json:"enter_on"`
}

// UserQueue is a queue of user
type UserQueue struct {
	Queue []User `json:"queue"`
}

// InsertNewUser insert a new user on the last position
func (q *UserQueue) InsertNewUser(user User) {
	q.Queue = append(q.Queue, user)
}

// Remove removes the first user on the top of the list
func (q *UserQueue) Remove() User {
	t := q.Queue
	q.Queue = t[1:]

	return t[0]
}

// Len returns the length of the list
func (q *UserQueue) Len() int {
	return len(q.Queue)
}

// GetPlaceByCPF returns the position of the user on the list or -1 if not found
func (q *UserQueue) GetPlaceByCPF(cpf string) int {
	for i, u := range q.Queue {
		if u.Cpf == cpf {
			return i + 1
		}
	}

	return -1
}
