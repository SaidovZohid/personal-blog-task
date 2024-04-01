package models

type Comment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	PostID    int64     `json:"post_id"`
	CreatedAt string    `json:"created_at"`
	Info      *UserData `json:"user_info"`
}

type UpdateComment struct {
	Content string `json:"content"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
	PostId  int64  `json:"post_id"`
}

type GetAllCommentsResponse struct {
	Comments []*Comment `json:"comments"`
	Count    int64      `json:"count"`
}

type GetAllCommentsParams struct {
	PostId int64 `json:"post_id" binding:"required"`
}
