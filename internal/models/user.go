package models

import "time"

type User struct {
	ID   int64     `json:"id"`
	Name string    `json:"name" validate:"required,min=1,max=200"`
	Dob  time.Time `json:"dob" validate:"required"`
}

type UserResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
	Age  int    `json:"age,omitempty"`
}
