package repo

import (
	"context"
	"time"
)

type Comment struct {
	ID        int64
	PostID    int64
	UserID    int64
	Content   string
	CreatedAt time.Time
	UserInfo  *UserInfo
}

type UpdateComment struct {
	ID      int64
	UserID  int64
	Content string
}

type CommentStorageI interface {
	Create(ctx context.Context, u *Comment) (*Comment, error)
	Update(ctx context.Context, u *UpdateComment) error
	Delete(ctx context.Context, commentId int64, userId int64) error
	GetAll(ctx context.Context, postId int64) (*GetAllCommentsResult, error)
}

type GetAllCommentsResult struct {
	Comments []*Comment
	Count    int64
}
