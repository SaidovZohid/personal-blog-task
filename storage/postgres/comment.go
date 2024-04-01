package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SaidovZohid/personal-blog-task/storage/repo"
	"github.com/jmoiron/sqlx"
)

type commentRepo struct {
	db *sqlx.DB
}

func NewCommentStorage(db *sqlx.DB) repo.CommentStorageI {
	return &commentRepo{
		db: db,
	}
}

func (c *commentRepo) Create(ctx context.Context, comment *repo.Comment) (*repo.Comment, error) {
	query := `
		INSERT INTO comments (
			post_id,
			user_id,
			content
		) VALUES ($1, $2, $3)
		RETURNING 
		id, 
		created_at
	`

	err := c.db.QueryRow(
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (c *commentRepo) Update(ctx context.Context, comment *repo.UpdateComment) error {
	query := `
		UPDATE comments SET
			content = $1
	    WHERE id = $2 AND user_id = $3
	`

	res, err := c.db.Exec(
		query,
		comment.Content,
		comment.ID,
		comment.UserID,
	)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (cr *commentRepo) Delete(ctx context.Context, commentId int64, userId int64) error {
	query := `
		DELETE FROM comments WHERE id = $1 AND user_id = $2
	`

	res, err := cr.db.Exec(
		query,
		commentId,
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

func (c *commentRepo) GetAll(ctx context.Context, postId int64) (*repo.GetAllCommentsResult, error) {
	result := repo.GetAllCommentsResult{
		Comments: make([]*repo.Comment, 0),
	}

	query := `
		SELECT
			c.id,
			c.post_id,
			c.user_id,
			c.content,
			c.created_at,
			u.name,
			u.id
		FROM comments c 
		INNER JOIN users u 	ON c.user_id = u.id  ` +
		fmt.Sprintf(" WHERE c.post_id = %d ", postId) + ` ORDER BY c.created_at DESC`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment repo.Comment
		var info repo.UserInfo
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&info.Name,
			&info.Id,
		)
		if err != nil {
			return nil, err
		}

		comment.UserInfo = &info
		result.Comments = append(result.Comments, &comment)
	}

	queryCount := "SELECT count(id) FROM comments " + fmt.Sprintf(" WHERE post_id = %d ", postId)

	err = c.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
