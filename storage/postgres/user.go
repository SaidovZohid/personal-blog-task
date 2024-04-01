package postgres

import (
	"context"
	"database/sql"

	"github.com/SaidovZohid/personal-blog-task/storage/repo"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

func (u userRepo) Create(ctx context.Context, user *repo.User) (*repo.User, error) {
	query := `
		INSERT INTO users (
			name,
			email,
			password,
			role
		) VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at
	`
	err := u.db.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
	).Scan(
		&user.Id,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u userRepo) Get(ctx context.Context, userId int64) (*repo.User, error) {
	var user repo.User
	query := `
		SELECT
			id,
			name,
			email,
			password,
			role,
			created_at
		FROM users WHERE id = $1	
	`
	err := u.db.QueryRow(
		query,
		userId,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u userRepo) GetByEmail(ctx context.Context, email string) (*repo.User, error) {
	var user repo.User
	query := `
		SELECT
			id,
			name,
			email,
			password,
			role,
			created_at
		FROM users WHERE email = $1	
	`
	err := u.db.QueryRow(
		query,
		email,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u userRepo) Update(ctx context.Context, user *repo.UpdateUser) (*repo.User, error) {
	var res repo.User
	query := `
		UPDATE users SET 
			name = $1
		WHERE id = $2 
		RETURNING id, name, email, role, created_at
	`
	err := u.db.QueryRow(
		query,
		user.Name,
		user.Id,
	).Scan(
		&res.Id,
		&res.Name,
		&res.Email,
		&res.Role,
		&res.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (u userRepo) Delete(ctx context.Context, userId int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := u.db.Exec(
		query,
		userId,
	)
	if err != nil {
		return err
	}
	if res, _ := result.RowsAffected(); res == 0 {
		return sql.ErrNoRows
	}

	return nil
}
