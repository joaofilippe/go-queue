package models

import "time"

// User is the base model
type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Cpf      string   `json:"cpf"`
	Phone    string   `json:"phone"`
	Priority bool     `json:"priority"`
	Symptoms []string `json:"symptoms"`
	Others   string   `json:"others"`
	Next     []User `json:"next"`

	Place   int       `json:"level"`
	EnterOn time.Time `json:"enter_on"`
}

// UserForm is the form model dto
type UserForm struct {
	Name  string `form:"name" json:"name"`
	Cpf   string `form:"cpf" json:"cpf"`
	Phone string `form:"phone" json:"phone"`

	Priority     string `form:"priority" json:"priority"`
	Fever        string `form:"Fever" json:"Fever"`
	Nausea       string `form:"nausea" json:"nausea"`
	Headache     string `form:"headache" json:"headache"`
	HeartDisease string `form:"heart_disease" json:"heart_disease"`
	Air          string `form:"air" json:"air"`
	Throat       string `form:"throat" json:"throat"`
	Hypertension string `form:"hypertension" json:"hypertension"`
	LowPressure  string `form:"low_pressure" json:"low_pressure"`
	Vertigo      string `form:"vertigo" json:"vertigo"`
	Alergy       string `form:"alergy" json:"alergy"`
	Diabetes     string `form:"diabetes" json:"diabetes"`
	Others       string `form:"others" json:"others"`
}
