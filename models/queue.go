package models

import "time"

// User is the base model
type User struct {
	ID           int       `form:"id" json:"id"`
	Name         string    `form:"name" json:"name"`
	Cpf          string    `form:"cpf" json:"cpf"`
	Phone        string    `form:"phone" json:"phone"`
	SymptomsList []string  `form:"symtoms_list" json:"symptoms_list"`
	Description  string    `form:"others" json:"others"`
	Level        int       `form:"level" json:"level"`
	EnterOn      time.Time `form:"enter_on" json:"enter_on"`

	Priority      string `form:"priority" json:"priority"`
	Temperature   string `form:"temperature" json:"temperature"`
	Nausea        string `form:"nausea" json:"nausea"`
	Heart         string `form:"heart" json:"heart"`
	Air           string `form:"air" json:"air"`
	Throat        string `form:"throat" json:"throat"`
	Headache      string `form:"headache" json:"headache"`
	Hypertension  string `form:"hypertension" json:"hypertension"`
	Diabetes      string `form:"diabetes" json:"diabetes"`
	Obesity       string `form:"obesity" json:"obesity"`
	BloodPressure string `form:"blood_pressure" json:"blood_pressure"`
	Vertigo       string `form:"vertigo" json:"vertigo"`
	HeartDisease  string `form:"heart-disease" json:"heart-disease"`
	Alergy        string `form:"alergy" json:"alergy"`
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

// GetPlaceByID returns the position of the user on the list or -1 if not found
func (q *UserQueue) GetPlaceByID(id int) int {
	for i, u := range q.Queue {
		if u.ID == id {
			return i + 1
		}
	}

	return -1
}
