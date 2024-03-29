package repo

import (
	"context"
	"time"
)

type UserStorageI interface {
	Create(ctx context.Context, u *User) (*User, error)
	Get(ctx context.Context, userId int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *UpdateUser) (*User, error)
	Delete(ctx context.Context, userId int) error
}

type User struct {
	Id        int
	Name      string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
}

type UpdateUser struct {
	Id   int
	Name string
}
