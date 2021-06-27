package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	Create(ctx context.Context, randomCode, url string) (string, error)
	FindURLByCode(ctx context.Context, code string) (string, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return repo{
		db,
	}
}

func (r repo) Create(ctx context.Context, randomCode, url string) (string, error) {
	query := `INSERT INTO
		urls (code, url)
		VALUES ($1, $2)
		RETURNING code`
	fmt.Println("randomCode", randomCode)
	_, err := r.db.ExecContext(ctx, query, randomCode, url)

	return randomCode, err
}

func (r repo) FindURLByCode(ctx context.Context, code string) (string, error) {
	query := `SELECT
		url
		FROM urls
		WHERE code = $1
		ORDER BY id DESC LIMIT 1`

	var url string

	err := r.db.GetContext(ctx, &url, query, code)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", ErrInvalidCode
	case err != nil:
		return "", err
	default:
		return url, nil
	}
}
