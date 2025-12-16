package repository

import (
	"context"
	"database/sql"
	"time"
)

type UserRow struct {
	ID   int64
	Name string
	Dob  time.Time
}

type Queries interface {
	CreateUser(ctx context.Context, name string, dob time.Time) (UserRow, error)
	GetUser(ctx context.Context, id int64) (UserRow, error)
	UpdateUser(ctx context.Context, id int64, name string, dob time.Time) (UserRow, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context) ([]UserRow, error)
}

// Fallback simple implementation until sqlc is generated.
// This lets the app run; swap with sqlc-generated code when ready.

type SimpleRepo struct{ db *sql.DB }

func NewSimpleRepo(db *sql.DB) *SimpleRepo { return &SimpleRepo{db: db} }

func (r *SimpleRepo) CreateUser(ctx context.Context, name string, dob time.Time) (UserRow, error) {
	var u UserRow
	err := r.db.QueryRowContext(ctx, "INSERT INTO users(name, dob) VALUES($1,$2) RETURNING id, name, dob", name, dob).Scan(&u.ID, &u.Name, &u.Dob)
	return u, err
}

func (r *SimpleRepo) GetUser(ctx context.Context, id int64) (UserRow, error) {
	var u UserRow
	err := r.db.QueryRowContext(ctx, "SELECT id, name, dob FROM users WHERE id=$1", id).Scan(&u.ID, &u.Name, &u.Dob)
	return u, err
}

func (r *SimpleRepo) UpdateUser(ctx context.Context, id int64, name string, dob time.Time) (UserRow, error) {
	var u UserRow
	err := r.db.QueryRowContext(ctx, "UPDATE users SET name=$2, dob=$3 WHERE id=$1 RETURNING id, name, dob", id, name, dob).Scan(&u.ID, &u.Name, &u.Dob)
	return u, err
}

func (r *SimpleRepo) DeleteUser(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	return err
}

func (r *SimpleRepo) ListUsers(ctx context.Context) ([]UserRow, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, dob FROM users ORDER BY id")
	if err != nil { return nil, err }
	defer rows.Close()
	var res []UserRow
	for rows.Next() {
		var u UserRow
		if err := rows.Scan(&u.ID, &u.Name, &u.Dob); err != nil { return nil, err }
		res = append(res, u)
	}
	return res, rows.Err()
}
