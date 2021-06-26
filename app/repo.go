package app

import "context"

type Repo interface {
	FindOrCreateCode(ctx context.Context, randomCode, url string) (string, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo{
	return repo{
		db
	}
}

func (r repo) FindOrCreateCode(ctx context.Context, randomCode, url string) (string, error){
	query := `WITH new_url AS (
		INSERT INTO urls (code, url)
		VALUES ($1, $2)
		ON CONFLICT (url) DO NOTHING
		RETURNING code
	) SELECT COALESCE(
		(SELECT code FROM new_url),
		(SELECT code FROM urls WHERE url = $2)
	)`

	var code string

	err := r.db.GetContext(ctx, &code, query, randomCode, url)

	return code, err
}