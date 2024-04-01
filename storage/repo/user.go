package repo

import (
	"context"
	"time"
)

type UserStorageI interface {
	Create(ctx context.Context, u *User) (*User, error)
	Get(ctx context.Context, userId int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *UpdateUser) (*User, error)
	Delete(ctx context.Context, userId int) error
}

type User struct {
	Id        int64
	Name      string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
}

type UserInfo struct {
	Id   int64
	Name string
}

type UpdateUser struct {
	Id   int64
	Name string
}
