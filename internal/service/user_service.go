package service

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	db "github.com/rdssj/golang-assessment/db/sqlc"
	"github.com/rdssj/golang-assessment/internal/models"
)

type UserService struct {
	repo     db.Querier
	validate *validator.Validate
	clockNow func() time.Time
}

func NewUserService(repo db.Querier) *UserService {
	return &UserService{repo: repo, validate: validator.New(), clockNow: time.Now}
}

func (s *UserService) Create(ctx context.Context, m models.User) (models.User, error) {
	if err := s.validate.Struct(m); err != nil {
		return models.User{}, err
	}
	u, err := s.repo.CreateUser(ctx, db.CreateUserParams{Name: m.Name, Dob: m.Dob})
	if err != nil {
		return models.User{}, err
	}
	return models.User{ID: u.ID, Name: u.Name, Dob: u.Dob}, nil
}

func (s *UserService) Get(ctx context.Context, id int64) (models.UserResponse, error) {
	u, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return models.UserResponse{}, err
	}
	now := s.clockNow()
	age := CalculateAge(u.Dob, now)
	return models.UserResponse{ID: u.ID, Name: u.Name, Dob: u.Dob.Format("2006-01-02"), Age: age}, nil
}

func (s *UserService) Update(ctx context.Context, id int64, m models.User) (models.User, error) {
	if err := s.validate.Struct(m); err != nil {
		return models.User{}, err
	}
	u, err := s.repo.UpdateUser(ctx, db.UpdateUserParams{ID: id, Name: m.Name, Dob: m.Dob})
	if err != nil {
		return models.User{}, err
	}
	return models.User{ID: u.ID, Name: u.Name, Dob: u.Dob}, nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) List(ctx context.Context, limit, offset int32) ([]models.UserResponse, error) {
	rows, err := s.repo.ListUsers(ctx, db.ListUsersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	now := s.clockNow()
	out := make([]models.UserResponse, 0, len(rows))
	for _, u := range rows {
		out = append(out, models.UserResponse{ID: u.ID, Name: u.Name, Dob: u.Dob.Format("2006-01-02"), Age: CalculateAge(u.Dob, now)})
	}
	return out, nil
}
