package repository

import (
	"alumnihub/internal/models"
	"context"
	"time"
)

func (m *PostgresDBRepo) InsertPasswordReset(pr models.PasswordReset) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into password_resets (email, token, expires_at, created_at)
			values ($1, $2, $3, $4)`

	_, err := m.DB.ExecContext(ctx, stmt,
		pr.Email,
		pr.Token,
		pr.ExpiresAt,
		pr.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetPasswordResetByToken(token string) (*models.PasswordReset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, email, token, expires_at, created_at 
			from password_resets where token = $1`

	var pr models.PasswordReset
	row := m.DB.QueryRowContext(ctx, query, token)

	err := row.Scan(
		&pr.ID,
		&pr.Email,
		&pr.Token,
		&pr.ExpiresAt,
		&pr.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &pr, nil
}

func (m *PostgresDBRepo) DeletePasswordResetByToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from password_resets where token = $1`

	_, err := m.DB.ExecContext(ctx, stmt, token)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeletePasswordResetsByEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from password_resets where email = $1`

	_, err := m.DB.ExecContext(ctx, stmt, email)
	if err != nil {
		return err
	}

	return nil
}

// Helper to get current time (used by service)
func (m *PostgresDBRepo) Now() time.Time {
	return time.Now()
}