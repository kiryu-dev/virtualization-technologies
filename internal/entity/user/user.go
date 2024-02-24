package user

import (
	"context"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Repository interface {
	GetAll(ctx context.Context, offset uint64, count uint64) ([]User, error)
	Get(ctx context.Context, id int) (*User, error)
	Create(ctx context.Context, user User) (int, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id int) (*User, error)
}
