package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SaidovZohid/personal-blog-task/storage/repo"
	"github.com/jmoiron/sqlx"
)

type replyRepo struct {
	db *sqlx.DB
}

func NewReplyStorage(db *sqlx.DB) repo.ReplyStorageI {
	return &replyRepo{
		db: db,
	}
}

func (r replyRepo) Create(ctx context.Context, reply *repo.Reply) (*repo.Reply, error) {
	query := `
		INSERT INTO replies (
			post_id,
			comment_id,
			user_id,
			content
		) VALUES ($1, $2, $3, $4)
		RETURNING 
		id, 
		created_at
	`

	err := r.db.QueryRow(
		query,
		reply.PostId,
		reply.CommentId,
		reply.UserId,
		reply.Content,
	).Scan(
		&reply.ID,
		&reply.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (r replyRepo) Update(ctx context.Context, updateR *repo.UpdateReply) error {
	query := `
		UPDATE replies SET
			content = $1
	    WHERE id = $2 AND user_id = $3
	`

	res, err := r.db.Exec(
		query,
		updateR.Content,
		updateR.Id,
		updateR.UserId,
	)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r replyRepo) Delete(ctx context.Context, replyId int64, userId int64) error {
	query := `
		DELETE FROM replies WHERE id = $1 AND user_id = $2
	`

	res, err := r.db.Exec(
		query,
		replyId,
		userId,
	)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r replyRepo) GetAll(ctx context.Context, reqParams *repo.GetAllRepliesParams) (*repo.GetAllRepliesResult, error) {
	result := repo.GetAllRepliesResult{
		Replies: make([]*repo.Reply, 0),
	}

	orderBy := ` ORDER BY r.created_at DESC `

	filter := " WHERE TRUE "
	if reqParams.CommentId > 0 {
		filter += fmt.Sprintf(" AND r.comment_id = %d ", reqParams.CommentId)
	}
	if reqParams.PostId > 0 {
		filter += fmt.Sprintf(" AND r.post_id = %d ", reqParams.PostId)
	}

	query := `
		SELECT
			r.id,
			r.comment_id,
			r.user_id,
			r.post_id,
			r.content,
			r.created_at,
			u.name,
			u.id
		FROM replies r
		INNER JOIN users u 	ON r.user_id = u.id  ` +
		filter + orderBy

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var res repo.Reply
		var info repo.UserInfo
		err := rows.Scan(
			&res.ID,
			&res.CommentId,
			&res.UserId,
			&res.PostId,
			&res.Content,
			&res.CreatedAt,
			&info.Name,
			&info.Id,
		)
		if err != nil {
			return nil, err
		}
		res.UserInfo = &info
		result.Replies = append(result.Replies, &res)
	}

	queryCount := "SELECT count(id) FROM replies r " + filter

	err = r.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
