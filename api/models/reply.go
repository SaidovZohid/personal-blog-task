package models

type Reply struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	PostId    int64     `json:"post_id"`
	CommentId int64     `json:"comment_id"`
	CreatedAt string    `json:"created_at"`
	Info      *UserData `json:"user_info"`
}

type UpdateReply struct {
	Content string `json:"content"`
}

type CreateReplyRequest struct {
	Content   string `json:"content"`
	CommentId int64  `json:"comment_id"`
	PostId    int64  `json:"post_id"`
}

type GetAllRepliesResponse struct {
	Replies []*Reply `json:"replies"`
	Count   int64    `json:"count"`
}

type GetAllRepliesParams struct {
	CommentId int64 `json:"comment_id" binding:"required"`
	PostId    int64 `json:"post_id"`
}
