package repository

import (
    "alumnihub/internal/models"
    "context"
)
func (m *PostgresDBRepo) AllJobs() ([]*models.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, user_id, job_position, company, job_location, job_type, min_salary, max_salary, closed, description, created_at, updated_at from jobs order by id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job

	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID,
			&job.UserID,
			&job.JobPosition,
			&job.Company,
			&job.JobLocation,
			&job.JobType,
			&job.MinSalary,
			&job.MaxSalary,
			&job.Closed,
			&job.Description,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}


func (m *PostgresDBRepo) Job(id int) (*models.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, user_id, job_position, company, job_location, job_type, min_salary, max_salary, closed, description, created_at, updated_at from jobs where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var job models.Job

	err := row.Scan(
		&job.ID,
		&job.UserID,
		&job.JobPosition,
		&job.Company,
		&job.JobLocation,
		&job.JobType,
		&job.MinSalary,
		&job.MaxSalary,
		&job.Closed,
		&job.Description,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &job, nil
}


func (m *PostgresDBRepo) InsertJob(job models.Job) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `insert into jobs (user_id, job_position, company, job_location, job_type, min_salary, max_salary, description, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		job.UserID,
		job.JobPosition,
		job.Company,
		job.JobLocation,
		job.JobType,
		job.MinSalary,
		job.MaxSalary,
		job.Description,
		job.CreatedAt,
		job.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}


func (m *PostgresDBRepo) UpdateJob(job models.Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `update jobs set job_position = $1, company = $2, job_location = $3, job_type = $4, min_salary = $5, max_salary = $6, description = $7, closed = $8, updated_at = $9 
			where id = $10`

	_, err := m.DB.ExecContext(ctx, stmt,
		job.JobPosition,
		job.Company,
		job.JobLocation,
		job.JobType,
		job.MinSalary,
		job.MaxSalary,
		job.Description,
		job.Closed,
		job.UpdatedAt,
		job.ID,
	)

	if err != nil {
		return err
	}

	return nil
}


func (m *PostgresDBRepo) DeleteJob(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from jobs where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}