package repository

import (
	"context"
	"database/sql"
)

func (m *PostgresDBRepo) CountAlumni() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `SELECT COUNT(id) as count_alumni FROM alumni`

	row := m.DB.QueryRowContext(ctx, query)

	var count int

	err := row.Scan(
		&count,
	)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if err == sql.ErrNoRows {
		return 0, nil
	}

	return count, nil

}

func (m *PostgresDBRepo) CountAlumniAccount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `SELECT COUNT(id) as count_alumni FROM alumni_profile`

	row := m.DB.QueryRowContext(ctx, query)

	var count int

	err := row.Scan(
		&count,
	)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if err == sql.ErrNoRows {
		return 0, nil
	}

	return count, nil

}
