package repository

import (
	"alumnihub/internal/models"
	"context"
)

func (m *PostgresDBRepo) AllArticles() ([]*models.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, title, slug, body, image, status, created_at, updated_at, published_at from articles order by id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*models.Article

	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Slug,
			&article.Body,
			&article.Image,
			&article.Status,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.PublishedAt,
		)
		if err != nil {
			return nil, err
		}

		articles = append(articles, &article)
	}

	return articles, nil
}

func (m *PostgresDBRepo) Article(id int) (*models.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, title, slug, body, image, status, created_at, updated_at, published_at
				FROM articles
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var article models.Article

	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Body,
		&article.Image,
		&article.Status,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.PublishedAt,
	)

	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (m *PostgresDBRepo) ArticleBySlug(slug string) (*models.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, title, slug, body, image, status, created_at, updated_at, published_at
				FROM articles
				WHERE slug = $1
			`

	row := m.DB.QueryRowContext(ctx, query, slug)

	var article models.Article

	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Body,
		&article.Image,
		&article.Status,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.PublishedAt,
	)

	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (m *PostgresDBRepo) InsertArticle(article models.Article) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into articles (title, slug, body, image, status, created_at, updated_at, published_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		article.Title,
		article.Slug,
		article.Body,
		article.Image,
		article.Status,
		article.CreatedAt,
		article.UpdatedAt,
		article.PublishedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) UpdateArticle(article models.Article) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `update articles set title = $1, slug = $2, body = $3,
				status = $4, updated_at = $5, published_at = $6, image = $7 
				where id = $8`

	_, err := m.DB.ExecContext(ctx, stmt,
		article.Title,
		article.Slug,
		article.Body,
		article.Status,
		article.UpdatedAt,
		article.PublishedAt,
		article.Image,
		article.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteArticle(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from articles where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
