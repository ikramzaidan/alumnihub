package repository

import (
	"alumnihub/internal/models"
	"context"
)

func (m *PostgresDBRepo) AllAlumni() ([]*models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, nisn, nis, name, gender, phone, graduation_year, class from alumni order by id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumnis []*models.Alumni

	for rows.Next() {
		var alumni models.Alumni
		err := rows.Scan(
			&alumni.ID,
			&alumni.NISN,
			&alumni.NIS,
			&alumni.Name,
			&alumni.Gender,
			&alumni.Phone,
			&alumni.Year,
			&alumni.Class,
		)
		if err != nil {
			return nil, err
		}

		alumnis = append(alumnis, &alumni)
	}

	return alumnis, nil
}

func (m *PostgresDBRepo) Alumni(id int) (*models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `
				SELECT id, nisn, nis, name, gender, phone, graduation_year, class
				FROM alumni
				WHERE id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, id)

	var alumni models.Alumni

	err := row.Scan(
		&alumni.ID,
		&alumni.NISN,
		&alumni.NIS,
		&alumni.Name,
		&alumni.Gender,
		&alumni.Phone,
		&alumni.Year,
		&alumni.Class,
	)

	if err != nil {
		return nil, err
	}

	return &alumni, nil
}

func (m *PostgresDBRepo) InsertAlumni(alumni models.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into alumni (nisn, nis, name, gender, phone, graduation_year, class)
			values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		alumni.NISN,
		alumni.NIS,
		alumni.Name,
		alumni.Gender,
		alumni.Phone,
		alumni.Year,
		alumni.Class,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) UpdateAlumni(alumni models.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `update alumni set name = $1, gender = $2, phone = $3, graduation_year = $4, class = $5, nisn = $6, nis = $7 
			where id = $8`

	_, err := m.DB.ExecContext(ctx, stmt,
		alumni.Name,
		alumni.Gender,
		alumni.Phone,
		alumni.Year,
		alumni.Class,
		alumni.NISN,
		alumni.NIS,
		alumni.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteAlumni(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from alumni where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetAlumniByNISN(nisn string) (*models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, nisn, nis, name, gender, phone, graduation_year, class 
			from alumni where nisn = $1`

	var alumni models.Alumni
	row := m.DB.QueryRowContext(ctx, query, nisn)

	err := row.Scan(
		&alumni.ID,
		&alumni.NISN,
		&alumni.NIS,
		&alumni.Name,
		&alumni.Gender,
		&alumni.Phone,
		&alumni.Year,
		&alumni.Class,
	)

	if err != nil {
		return nil, err
	}

	return &alumni, nil

}

func (m *PostgresDBRepo) GetAlumniNameByID(id int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select name from alumni where id = $1`

	var alumni models.Alumni
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&alumni.Name,
	)

	if err != nil {
		return "", err
	}

	return alumni.Name, nil
}
