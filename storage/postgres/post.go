package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SaidovZohid/personal-blog-task/storage/repo"
	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPostStorage(db *sqlx.DB) repo.PostStorageI {
	return &postRepo{
		db: db,
	}
}

func (p postRepo) Create(ctx context.Context, post *repo.Post) (*repo.Post, error) {
	query := `
		INSERT INTO posts (
			header,
			content,
			user_id
		) VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`
	err := p.db.QueryRow(
		query,
		post.Header,
		post.Body,
		post.UserId,
	).Scan(
		&post.Id,
		&post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (p postRepo) Get(ctx context.Context, postId int64) (*repo.Post, error) {
	var post repo.Post
	query := `
		SELECT
			p.id,
			p.header,
			p.content,
			p.user_id,
			p.created_at,
			u.name,
			u.id
		FROM posts p INNER JOIN users u	ON p.user_id = u.id WHERE p.id = $1
	`
	var userInfo repo.UserInfo
	err := p.db.QueryRow(
		query,
		postId,
	).Scan(
		&post.Id,
		&post.Header,
		&post.Body,
		&post.UserId,
		&post.CreatedAt,
		&userInfo.Name,
		&userInfo.Id,
	)
	if err != nil {
		return nil, err
	}
	post.UserInfo = &userInfo

	return &post, nil
}

func (p postRepo) Update(ctx context.Context, post *repo.UpdatePost) (*repo.Post, error) {
	var res repo.Post
	query := `
		UPDATE posts SET 
			header = $1,
			content = $2
		WHERE id = $3 AND user_id = $4
		RETURNING id, header, content, user_id, created_at
	`
	err := p.db.QueryRow(
		query,
		post.Header,
		post.Body,
		post.Id,
		post.UserId,
	).Scan(
		&res.Id,
		&res.Header,
		&res.Body,
		&res.UserId,
		&res.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (p postRepo) Delete(ctx context.Context, postId int64, userID int64) error {
	query := `DELETE FROM posts WHERE id = $1 AND user_id = $2`
	result, err := p.db.Exec(
		query,
		postId,
		userID,
	)
	if err != nil {
		return err
	}
	if res, _ := result.RowsAffected(); res == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p postRepo) GetAll(ctx context.Context, params *repo.GetAllParams) (*repo.GetAllResult, error) {
	result := repo.GetAllResult{
		Result: make([]*repo.Post, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, offset)

	filter := " WHERE true"

	if params.UserId != -1 {
		filter += fmt.Sprintf(" AND p.user_id = %d", params.UserId)
	}

	orderBy := " ORDER BY created_at DESC"
	if params.SortByDate != "" {
		orderBy = fmt.Sprintf(" ORDER BY created_at %s", params.SortByDate)
	}

	query := `
		SELECT
			p.id,
			p.header,
			p.content,
			p.user_id,
			p.created_at,
			u.name,
			u.id
		FROM posts p INNER JOIN users u ON p.user_id = u.id
	` + filter + orderBy + limit

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post repo.Post
		var userInfo repo.UserInfo
		err := rows.Scan(
			&post.Id,
			&post.Header,
			&post.Body,
			&post.UserId,
			&post.CreatedAt,
			&userInfo.Name,
			&userInfo.Id,
		)
		if err != nil {
			return nil, err
		}

		post.UserInfo = &userInfo
		result.Result = append(result.Result, &post)
	}

	queryCount := "SELECT count(id) FROM posts p " + filter

	err = p.db.QueryRow(queryCount).Scan(&result.Count)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
