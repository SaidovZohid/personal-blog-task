package storage

import (
	"github.com/SaidovZohid/personal-blog-task/storage/postgres"
	"github.com/SaidovZohid/personal-blog-task/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
	Post() repo.PostStorageI
	Comment() repo.CommentStorageI
	Reply() repo.ReplyStorageI
}

type storage struct {
	userRepo    repo.UserStorageI
	postRepo    repo.PostStorageI
	commentRepo repo.CommentStorageI
	replyRepo   repo.ReplyStorageI
}

func NewStorage(db *sqlx.DB) StorageI {
	return &storage{
		userRepo:    postgres.NewUserStorage(db),
		postRepo:    postgres.NewPostStorage(db),
		commentRepo: postgres.NewCommentStorage(db),
		replyRepo:   postgres.NewReplyStorage(db),
	}
}

func (s storage) User() repo.UserStorageI {
	return s.userRepo
}

func (s storage) Post() repo.PostStorageI {
	return s.postRepo
}

func (s storage) Comment() repo.CommentStorageI {
	return s.commentRepo
}

func (s storage) Reply() repo.ReplyStorageI {
	return s.replyRepo
}
