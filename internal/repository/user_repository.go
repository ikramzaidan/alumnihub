package repository

import (
	"alumnihub/internal/models"
	"context"
	"database/sql"
)

func (m *PostgresDBRepo) InsertUser(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into users (username, email, password, is_admin, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6) returning id`

	var userID int

	err := m.DB.QueryRowContext(ctx, stmt,
		user.Username,
		user.Email,
		user.Password,
		user.IsAdmin,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (m *PostgresDBRepo) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, username, email, password, is_admin, created_at, updated_at 
			from users where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, username, email, password, is_admin, created_at, updated_at 
			from users where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (m *PostgresDBRepo) GetUserUsernameByID(id int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select username from users where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.Username,
	)

	if err != nil {
		return "", err
	}

	return user.Username, nil
}

func (m *PostgresDBRepo) GetUserIDByUsername(username string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id from users where username = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, username)

	err := row.Scan(
		&user.ID,
	)

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (m *PostgresDBRepo) GetUserPhotoByID(id int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select photo from users where id = $1`

	var photo sql.NullString
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&photo)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	if photo.Valid {
		return photo.String, nil
	}

	return "", nil
}
