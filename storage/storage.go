package storage

import (
	"github.com/SaidovZohid/personal-blog-task/storage/postgres"
	"github.com/SaidovZohid/personal-blog-task/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
}

type storage struct {
	userRepo repo.UserStorageI
}

func NewStorage(db *sqlx.DB) StorageI {
	return &storage{
		userRepo: postgres.NewUserStorage(db),
	}
}

func (s storage) User() repo.UserStorageI {
	return s.userRepo
}
