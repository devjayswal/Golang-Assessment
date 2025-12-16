// Code generated manually to mimic sqlc interfaces.
package db

import (
	"context"
	"time"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUser(ctx context.Context, id int64) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
}

type CreateUserParams struct {
	Name string    `json:"name"`
	Dob  time.Time `json:"dob"`
}

type UpdateUserParams struct {
	ID   int64     `json:"id"`
	Name string    `json:"name"`
	Dob  time.Time `json:"dob"`
}

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}
