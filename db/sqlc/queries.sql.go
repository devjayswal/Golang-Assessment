// Code generated manually to approximate sqlc generated code.
package db

import (
	"context"
	"database/sql"
)

// Queries provides type-safe DB operations.
type Queries struct {
	db *sql.DB
}

func New(db *sql.DB) *Queries {
	return &Queries{db: db}
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, "INSERT INTO users (name, dob) VALUES ($1, $2) RETURNING id, name, dob", arg.Name, arg.Dob)
	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Dob)
	return u, err
}

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, "SELECT id, name, dob FROM users WHERE id = $1", id)
	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Dob)
	return u, err
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, "UPDATE users SET name = $2, dob = $3 WHERE id = $1 RETURNING id, name, dob", arg.ID, arg.Name, arg.Dob)
	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Dob)
	return u, err
}

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT id, name, dob FROM users ORDER BY id LIMIT $1 OFFSET $2", arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Dob); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
