package models

type Post struct {
	ID        int64     `json:"id"`
	Header    string    `json:"header"`
	Body      string    `json:"body"`
	UserID    int64     `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UserInfo  *UserData `json:"user_info"`
}

type CreatePostRequest struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}

type UpdatePostRequest struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}

type GetAllPostsParams struct {
	Limit      int64  `json:"limit" binding:"required" default:"10"`
	Page       int64  `json:"page" binding:"required" default:"1"`
	UserId     int64  `json:"user_id"`
	SortByDate string `json:"sort" enums:"desc,asc" default:"desc"`
}

type GetAllPostsResponse struct {
	Posts []*Post `json:"posts"`
	Count int64   `json:"count"`
}

type GetPostInfo struct {
	ID        int64         `json:"id"`
	Header    string        `json:"header"`
	Body      string        `json:"body"`
	UserID    int64         `json:"user_id"`
	CreatedAt string        `json:"created_at"`
	UserInfo  *UserData     `json:"user_info"`
	Comments  *PostComments `json:"all_comments"`
}

type PostComments struct {
	Comments []*CommentWithReplies `json:"comments"`
	Count    int64                 `json:"count"`
}

type CommentWithReplies struct {
	ID        int64                  `json:"id"`
	Content   string                 `json:"content"`
	UserID    int64                  `json:"user_id"`
	PostID    int64                  `json:"post_id"`
	CreatedAt string                 `json:"created_at"`
	Info      *UserData              `json:"user_info"`
	Replies   *GetAllRepliesResponse `json:"all_replies"`
}
