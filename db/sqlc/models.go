// Code generated manually to mirror sqlc output for this project.
package db

import "time"

type User struct {
	ID   int64     `json:"id"`
	Name string    `json:"name"`
	Dob  time.Time `json:"dob"`
}
