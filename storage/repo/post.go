package repo

import (
	"context"
	"time"
)

type PostStorageI interface {
	Create(ctx context.Context, u *Post) (*Post, error)
	Get(ctx context.Context, postId int64) (*Post, error)
	Update(ctx context.Context, u *UpdatePost) (*Post, error)
	Delete(ctx context.Context, postIf int64, userId int64) error
	GetAll(ctx context.Context, param *GetAllParams) (*GetAllResult, error)
}

type Post struct {
	Id        int64
	Header    string
	Body      string
	UserId    int64
	CreatedAt time.Time
	UserInfo  *UserInfo
}

type UpdatePost struct {
	Id     int64
	Header string
	Body   string
	UserId int64
}

type GetAllParams struct {
	Page       int64
	Limit      int64
	UserId     int64
	SortByDate string
}

type GetAllResult struct {
	Result []*Post
	Count  int64
}
