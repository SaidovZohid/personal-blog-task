package repo

import (
	"context"
	"time"
)

type ReplyStorageI interface {
	Create(ctx context.Context, u *Reply) (*Reply, error)
	Update(ctx context.Context, u *UpdateReply) error
	Delete(ctx context.Context, commentId int64, userId int64) error
	GetAll(ctx context.Context, params *GetAllRepliesParams) (*GetAllRepliesResult, error)
}

type Reply struct {
	ID        int64
	CommentId int64
	UserId    int64
	PostId    int64
	Content   string
	CreatedAt time.Time
	UserInfo  *UserInfo
}

type UpdateReply struct {
	Id      int64
	UserId  int64
	Content string
}

type GetAllRepliesParams struct {
	CommentId int64
	PostId    int64
}

type GetAllRepliesResult struct {
	Replies []*Reply
	Count   int64
}
